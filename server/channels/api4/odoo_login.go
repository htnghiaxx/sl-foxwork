// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package api4

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/shared/mlog"
	"github.com/mattermost/mattermost/server/v8/channels/utils"
)

type odooLoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRPCRespErr `json:"error,omitempty"`
}

type jsonRPCRespErr struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func loginOdoo(c *Context, w http.ResponseWriter, r *http.Request) {

	baseURL := strings.TrimRight(os.Getenv("MM_ODOO_BASE_URL"), "/")
	if baseURL == "" {
		c.Err = model.NewAppError("loginOdoo", "api.user.odoo_login.missing_base_url", nil, "", http.StatusInternalServerError)
		return
	}
	dbName := os.Getenv("MM_ODOO_DB")
	if dbName == "" {
		c.Err = model.NewAppError("loginOdoo", "api.user.odoo_login.missing_db", nil, "", http.StatusInternalServerError)
		return
	}
	jsonrpcPath := os.Getenv("MM_ODOO_JSONRPC_PATH")
	if jsonrpcPath == "" {
		jsonrpcPath = "/jsonrpc"
	}

	c.Logger.Debug("loginOdoo config", mlog.String("base_url", baseURL), mlog.String("db", dbName), mlog.String("jsonrpc_path", jsonrpcPath))
	// Feature flag qua ENV theo tài liệu odoo-login.md

	var req odooLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.SetInvalidParamWithErr("body", err)
		return
	}
	c.Logger.Debug("loginOdoo received payload", mlog.String("identifier", req.Identifier))
	if req.Identifier == "" || req.Password == "" {
		c.SetInvalidParam("identifier/password")
		return
	}
	c.Logger.Debug("loginOdoo received payload", mlog.String("identifier", req.Identifier))
	timeoutMs := 8000
	if v := os.Getenv("MM_ODOO_TIMEOUT_MS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			timeoutMs = n
		}
	}

	httpClient := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}

	// Step 1: authenticate to get uid (via /web/session/authenticate)
	authPayload := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"db":       dbName,
			"login":    req.Identifier,
			"password": req.Password,
		},
		ID: 1,
	}
	c.Logger.Debug("loginOdoo authenticate to get uid", mlog.String("url", baseURL+"/web/session/authenticate"))
	uid, appErr := callOdooAuthenticate(httpClient, baseURL+"/web/session/authenticate", authPayload)
	c.Logger.Debug("loginOdoo authenticate to get uid", mlog.Int("uid", uid))
	if appErr != nil {
		c.Logger.Error("loginOdoo authenticate upstream error", mlog.Err(appErr))
		c.Err = appErr
		return
	}
	if uid == 0 {
		c.Logger.Info("loginOdoo authentication failed", mlog.String("identifier", req.Identifier))
		c.Err = model.NewAppError("loginOdoo", "api.user.login.invalid_credentials_email_username", nil, "", http.StatusUnauthorized)
		return
	}
	c.Logger.Debug("loginOdoo authenticate success", mlog.Int("uid", uid))

	// Step 2: fetch user info
	email, username, displayName, appErr := fetchOdooUserInfo(httpClient, baseURL+jsonrpcPath, dbName, uid, req.Password)
	if appErr != nil {
		c.Logger.Error("loginOdoo fetch user info error", mlog.Err(appErr))
		c.Err = appErr
		return
	}
	c.Logger.Debug("loginOdoo user info", mlog.String("email", email), mlog.String("username", username))
	if username == "" {
		username = req.Identifier
	}
	email = fmt.Sprintf("odoo_%d@odoo.local", uid)
	// Find or create Mattermost user
	mmUser, _ := c.App.GetUserByEmail(email)
	created := false
	if mmUser == nil {
		c.Logger.Info("loginOdoo creating new user", mlog.String("email", email))
		// create new user
		u := &model.User{
			Username:    model.CleanUsername(c.Logger, strings.Split(username, "@")[0]),
			Email:       email,
			FirstName:   displayName,
			AuthService: "odoo",
		}
		if u.Username == "" {
			u.Username = model.NewId()[:12]
		}
		var err *model.AppError
		mmUser, err = c.App.CreateUser(c.AppContext, u)
		if err != nil {
			c.Logger.Error("loginOdoo create user error", mlog.Err(err))
			c.Err = err
			return
		}
		created = true
	}

	// Issue session
	isMobile := utils.IsMobileRequest(r)
	session, err := c.App.DoLogin(c.AppContext, w, r, mmUser, "", isMobile, false, false)
	if err != nil {
		c.Logger.Error("loginOdoo issue session error", mlog.Err(err))
		c.Err = err
		return
	}
	c.AppContext = c.AppContext.WithSession(session)

	// sanitize and respond
	mmUser.Sanitize(map[string]bool{})
	resp := map[string]any{
		"user_id":        mmUser.Id,
		"username":       mmUser.Username,
		"email":          mmUser.Email,
		"create":         created,
		"updated_fields": []string{},
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func callOdooAuthenticate(httpClient *http.Client, url string, payload jsonRPCRequest) (int, *model.AppError) {
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return 0, model.NewAppError("odooAuth", "api.user.odoo_login.upstream_error", nil, err.Error(), http.StatusBadGateway)
	}
	defer res.Body.Close()

	// If upstream returns 401/403, treat as invalid credentials to allow fallback
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return 0, nil
	}

	var rpcRes jsonRPCResponse
	if err := json.NewDecoder(res.Body).Decode(&rpcRes); err != nil {
		return 0, model.NewAppError("odooAuth", "api.user.odoo_login.decode_error", nil, err.Error(), http.StatusBadGateway)
	}
	if rpcRes.Error != nil {
		// Heuristic: treat AccessDenied and alike as invalid credentials (fallback to DB)
		msg := strings.ToLower(rpcRes.Error.Message)
		if strings.Contains(msg, "wrong") || strings.Contains(msg, "password") || strings.Contains(msg, "access denied") || strings.Contains(msg, "invalid") || strings.Contains(msg, "authentication failed") {
			return 0, nil
		}
		// Also inspect error.data
		if rpcRes.Error.Data != nil {
			if m, ok := rpcRes.Error.Data.(map[string]interface{}); ok {
				if v, ok := m["name"].(string); ok {
					if strings.Contains(strings.ToLower(v), "accessdenied") {
						return 0, nil
					}
				}
				if v, ok := m["message"].(string); ok {
					if strings.Contains(strings.ToLower(v), "access denied") || strings.Contains(strings.ToLower(v), "invalid") {
						return 0, nil
					}
				}
			}
		}
		return 0, model.NewAppError("odooAuth", "api.user.odoo_login.upstream_error", nil, rpcRes.Error.Message, http.StatusBadGateway)
	}
	if len(rpcRes.Result) == 0 || string(rpcRes.Result) == "null" {
		return 0, nil
	}
	// Expecting result to be an object containing uid field
	var resultObj map[string]interface{}
	if err := json.Unmarshal(rpcRes.Result, &resultObj); err != nil {
		return 0, model.NewAppError("odooAuth", "api.user.odoo_login.decode_uid_error", nil, err.Error(), http.StatusBadGateway)
	}
	if v, ok := resultObj["uid"]; ok {
		switch t := v.(type) {
		case float64:
			return int(t), nil
		case int:
			return t, nil
		}
	}
	return 0, nil
}

func fetchOdooUserInfo(httpClient *http.Client, url, db string, uid int, password string) (email, login, name string, appErr *model.AppError) {
	// execute_kw search_read on res.users
	args := []interface{}{db, uid, password, "res.users", "search_read", []interface{}{[]interface{}{[]interface{}{"id", "=", uid}}}, map[string]interface{}{"fields": []string{"id", "name", "login", "email", "partner_id"}}}
	payload := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "call",
		Params: map[string]interface{}{
			"service": "object",
			"method":  "execute_kw",
			"args":    args,
		},
		ID: 2,
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return "", "", "", model.NewAppError("odooUserInfo", "api.user.odoo_login.upstream_error", nil, err.Error(), http.StatusBadGateway)
	}
	defer res.Body.Close()

	var rpcRes jsonRPCResponse
	if err := json.NewDecoder(res.Body).Decode(&rpcRes); err != nil {
		return "", "", "", model.NewAppError("odooUserInfo", "api.user.odoo_login.decode_error", nil, err.Error(), http.StatusBadGateway)
	}
	if rpcRes.Error != nil {
		return "", "", "", model.NewAppError("odooUserInfo", "api.user.odoo_login.upstream_error", nil, rpcRes.Error.Message, http.StatusBadGateway)
	}
	if len(rpcRes.Result) == 0 {
		return "", "", "", model.NewAppError("odooUserInfo", "api.user.odoo_login.user_not_found", nil, "", http.StatusUnauthorized)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(rpcRes.Result, &rows); err != nil {
		return "", "", "", model.NewAppError("odooUserInfo", "api.user.odoo_login.decode_user_error", nil, err.Error(), http.StatusBadGateway)
	}
	if len(rows) == 0 {
		return "", "", "", model.NewAppError("odooUserInfo", "api.user.odoo_login.user_not_found", nil, "", http.StatusUnauthorized)
	}
	row := rows[0]
	getString := func(key string) string {
		if v, ok := row[key]; ok {
			if s, ok := v.(string); ok {
				return s
			}
		}
		return ""
	}
	name = getString("name")
	login = getString("login")
	email = getString("email")
	return
}
