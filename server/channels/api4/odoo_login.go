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

// Minimal shape we care about from Odoo authenticate result
type odooAuthResult struct {
	UID           int    `json:"uid"`
	IsSystem      bool   `json:"is_system"`
	IsAdmin       bool   `json:"is_admin"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	UserCompanies struct {
		AllowedCompanies map[string]struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"allowed_companies"`
	} `json:"user_companies"`
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
	authRes, appErr := callOdooAuthenticate(httpClient, baseURL+"/web/session/authenticate", authPayload)
	c.Logger.Debug("loginOdoo authenticate to get uid", mlog.Int("uid", authRes.UID))
	if appErr != nil {
		c.Logger.Error("loginOdoo authenticate upstream error", mlog.Err(appErr))
		c.Err = appErr
		return
	}
	if authRes.UID == 0 {
		c.Logger.Info("loginOdoo authentication failed", mlog.String("identifier", req.Identifier))
		c.Err = model.NewAppError("loginOdoo", "api.user.login.invalid_credentials_email_username", nil, "", http.StatusUnauthorized)
		return
	}
	c.Logger.Debug("loginOdoo authenticate success", mlog.Int("uid", authRes.UID), mlog.Bool("is_system", authRes.IsSystem), mlog.Bool("is_admin", authRes.IsAdmin))

	// Step 2: fetch user info
	email, username, displayName, appErr := fetchOdooUserInfo(httpClient, baseURL+jsonrpcPath, dbName, authRes.UID, req.Password)
	if appErr != nil {
		c.Logger.Error("loginOdoo fetch user info error", mlog.Err(appErr))
		c.Err = appErr
		return
	}
	c.Logger.Debug("loginOdoo user info", mlog.String("email", email), mlog.String("username", username))
	if username == "" {
		username = req.Identifier
	}
	email = fmt.Sprintf("odoo_%d@odoo.local", authRes.UID)
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

	// Step 3: sync Odoo companies to Mattermost teams and roles (using fields from authenticate result)
	// Build companies list from authRes.UserCompanies
	var authCompanies []companyMeta
	if len(authRes.UserCompanies.AllowedCompanies) > 0 {
		authCompanies = make([]companyMeta, 0, len(authRes.UserCompanies.AllowedCompanies))
		for _, v := range authRes.UserCompanies.AllowedCompanies {
			if v.ID != 0 && v.Name != "" {
				authCompanies = append(authCompanies, companyMeta{ID: v.ID, Name: v.Name})
			}
		}
	}
	c.Logger.Debug("loginOdoo companies from auth result", mlog.Int("count", len(authCompanies)))
	if len(authCompanies) > 0 {
		for _, comp := range authCompanies {
			teamName := slugify(comp.Name)
			if teamName == "" {
				continue
			}
			team, getErr := c.App.GetTeamByName(teamName)
			if getErr != nil {
				// create team if missing
				newTeam := &model.Team{
					Name:        teamName,
					DisplayName: comp.Name,
					Type:        model.TeamOpen,
				}
				team, getErr = c.App.CreateTeam(c.AppContext, newTeam)
				if getErr != nil {
					c.Logger.Warn("loginOdoo: create team failed", mlog.String("team", teamName), mlog.Err(getErr))
					continue
				}
			}
			// join member
			if _, _, jErr := c.App.AddUserToTeam(c.AppContext, team.Id, mmUser.Id, ""); jErr != nil {
				c.Logger.Warn("loginOdoo: add user to team failed", mlog.String("team", team.Name), mlog.Err(jErr))
				continue
			}
			// elevate team admin if both flags are true
			if authRes.IsSystem && authRes.IsAdmin {
				if _, rErr := c.App.UpdateTeamMemberRoles(c.AppContext, team.Id, mmUser.Id, "team_user team_admin"); rErr != nil {
					c.Logger.Warn("loginOdoo: set team admin failed", mlog.String("team", team.Name), mlog.Err(rErr))
				}
			}
		}
	}

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

func callOdooAuthenticate(httpClient *http.Client, url string, payload jsonRPCRequest) (odooAuthResult, *model.AppError) {
	var empty odooAuthResult
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return empty, model.NewAppError("odooAuth", "api.user.odoo_login.upstream_error", nil, err.Error(), http.StatusBadGateway)
	}
	defer res.Body.Close()

	// If upstream returns 401/403, treat as invalid credentials to allow fallback
	if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
		return empty, nil
	}

	var rpcRes jsonRPCResponse
	if err := json.NewDecoder(res.Body).Decode(&rpcRes); err != nil {
		return empty, model.NewAppError("odooAuth", "api.user.odoo_login.decode_error", nil, err.Error(), http.StatusBadGateway)
	}
	if rpcRes.Error != nil {
		// Heuristic: treat AccessDenied and alike as invalid credentials (fallback to DB)
		msg := strings.ToLower(rpcRes.Error.Message)
		if strings.Contains(msg, "wrong") || strings.Contains(msg, "password") || strings.Contains(msg, "access denied") || strings.Contains(msg, "invalid") || strings.Contains(msg, "authentication failed") {
			return empty, nil
		}
		// Also inspect error.data
		if rpcRes.Error.Data != nil {
			if m, ok := rpcRes.Error.Data.(map[string]interface{}); ok {
				if v, ok := m["name"].(string); ok {
					if strings.Contains(strings.ToLower(v), "accessdenied") {
						return empty, nil
					}
				}
				if v, ok := m["message"].(string); ok {
					if strings.Contains(strings.ToLower(v), "access denied") || strings.Contains(strings.ToLower(v), "invalid") {
						return empty, nil
					}
				}
			}
		}
		return empty, model.NewAppError("odooAuth", "api.user.odoo_login.upstream_error", nil, rpcRes.Error.Message, http.StatusBadGateway)
	}
	if len(rpcRes.Result) == 0 || string(rpcRes.Result) == "null" {
		return empty, nil
	}
	var out odooAuthResult
	if err := json.Unmarshal(rpcRes.Result, &out); err != nil {
		// Fallback parser for just uid
		var resultObj map[string]interface{}
		if err2 := json.Unmarshal(rpcRes.Result, &resultObj); err2 == nil {
			if v, ok := resultObj["uid"]; ok {
				switch t := v.(type) {
				case float64:
					out.UID = int(t)
				case int:
					out.UID = t
				}
			}
		} else {
			return empty, model.NewAppError("odooAuth", "api.user.odoo_login.decode_uid_error", nil, err.Error(), http.StatusBadGateway)
		}
	}
	return out, nil
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

// companyMeta represents minimal company info from Odoo
type companyMeta struct {
	ID   int
	Name string
}

// fetchOdooUserCompaniesAndRoles returns companies the user belongs to and admin flags if present in Odoo
func fetchOdooUserCompaniesAndRoles(httpClient *http.Client, url, db string, uid int, password string) ([]companyMeta, bool, bool, *model.AppError) {
	// First, read user fields: company_ids, is_system, is_admin (customs may be absent), groups_id
	args := []interface{}{db, uid, password, "res.users", "search_read", []interface{}{[]interface{}{[]interface{}{"id", "=", uid}}}, map[string]interface{}{"fields": []string{"id", "company_ids", "is_system", "is_admin", "groups_id"}}}
	payload := jsonRPCRequest{JSONRPC: "2.0", Method: "call", Params: map[string]interface{}{"service": "object", "method": "execute_kw", "args": args}, ID: 3}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := httpClient.Do(req)
	if err != nil {
		return nil, false, false, model.NewAppError("odooUserMeta", "api.user.odoo_login.upstream_error", nil, err.Error(), http.StatusBadGateway)
	}
	defer res.Body.Close()
	var rpcRes jsonRPCResponse
	if err := json.NewDecoder(res.Body).Decode(&rpcRes); err != nil {
		return nil, false, false, model.NewAppError("odooUserMeta", "api.user.odoo_login.decode_error", nil, err.Error(), http.StatusBadGateway)
	}
	if rpcRes.Error != nil {
		return nil, false, false, model.NewAppError("odooUserMeta", "api.user.odoo_login.upstream_error", nil, rpcRes.Error.Message, http.StatusBadGateway)
	}
	var rows []map[string]interface{}
	if err := json.Unmarshal(rpcRes.Result, &rows); err != nil || len(rows) == 0 {
		return nil, false, false, model.NewAppError("odooUserMeta", "api.user.odoo_login.decode_user_error", nil, "", http.StatusBadGateway)
	}
	row := rows[0]
	// Extract company_ids
	companyIDs := []int{}
	if v, ok := row["company_ids"]; ok {
		if arr, ok := v.([]interface{}); ok {
			for _, iid := range arr {
				switch t := iid.(type) {
				case float64:
					companyIDs = append(companyIDs, int(t))
				case int:
					companyIDs = append(companyIDs, t)
				}
			}
		}
	}
	// Extract flags (optional fields)
	isSystem := false
	isAdmin := false
	if v, ok := row["is_system"]; ok {
		if b, ok := v.(bool); ok {
			isSystem = b
		}
	}
	if v, ok := row["is_admin"]; ok {
		if b, ok := v.(bool); ok {
			isAdmin = b
		}
	}
	// Fallback: infer admin if groups_id contains any element whose display name contains "Settings" or "Administration" (best-effort)
	if !isAdmin {
		if v, ok := row["groups_id"]; ok {
			if arr, ok := v.([]interface{}); ok {
				// if user has many groups, we cannot resolve names without extra RPC; skip for brevity
				// keep default false
				_ = arr
			}
		}
	}

	if len(companyIDs) == 0 {
		return []companyMeta{}, isSystem, isAdmin, nil
	}
	// Fetch companies by ids
	companyArgs := []interface{}{db, uid, password, "res.company", "search_read", []interface{}{[]interface{}{[]interface{}{"id", "in", companyIDs}}}, map[string]interface{}{"fields": []string{"id", "name"}}}
	companyPayload := jsonRPCRequest{JSONRPC: "2.0", Method: "call", Params: map[string]interface{}{"service": "object", "method": "execute_kw", "args": companyArgs}, ID: 4}
	companyBody, _ := json.Marshal(companyPayload)
	creq, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(companyBody))
	creq.Header.Set("Content-Type", "application/json")
	cres, err := httpClient.Do(creq)
	if err != nil {
		return nil, isSystem, isAdmin, model.NewAppError("odooCompanies", "api.user.odoo_login.upstream_error", nil, err.Error(), http.StatusBadGateway)
	}
	defer cres.Body.Close()
	var compRPC jsonRPCResponse
	if err := json.NewDecoder(cres.Body).Decode(&compRPC); err != nil {
		return nil, isSystem, isAdmin, model.NewAppError("odooCompanies", "api.user.odoo_login.decode_error", nil, err.Error(), http.StatusBadGateway)
	}
	if compRPC.Error != nil {
		return nil, isSystem, isAdmin, model.NewAppError("odooCompanies", "api.user.odoo_login.upstream_error", nil, compRPC.Error.Message, http.StatusBadGateway)
	}
	var compRows []map[string]interface{}
	if err := json.Unmarshal(compRPC.Result, &compRows); err != nil {
		return nil, isSystem, isAdmin, model.NewAppError("odooCompanies", "api.user.odoo_login.decode_error", nil, err.Error(), http.StatusBadGateway)
	}
	companies := make([]companyMeta, 0, len(compRows))
	for _, r := range compRows {
		id := 0
		name := ""
		if v, ok := r["id"]; ok {
			switch t := v.(type) {
			case float64:
				id = int(t)
			case int:
				id = t
			}
		}
		if v, ok := r["name"]; ok {
			if s, ok := v.(string); ok {
				name = s
			}
		}
		if id != 0 && name != "" {
			companies = append(companies, companyMeta{ID: id, Name: name})
		}
	}
	return companies, isSystem, isAdmin, nil
}

// slugify converts company name to Mattermost team name (lowercase, alnum and dashes only)
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if s == "" {
		return ""
	}
	// replace spaces and underscores with dash
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.ReplaceAll(s, " ", "-")
	// remove invalid characters
	var b strings.Builder
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' {
			b.WriteRune(ch)
		}
	}
	res := b.String()
	res = strings.Trim(res, "-")
	if res == "" {
		return ""
	}
	// collapse multiple dashes
	for strings.Contains(res, "--") {
		res = strings.ReplaceAll(res, "--", "-")
	}
	return res
}
