## Kế hoạch tích hợp đăng nhập qua Odoo (SSO tối giản theo proxy backend)

Tài liệu này mô tả thay đổi cần thực hiện để cho phép Mattermost xác thực người dùng thông qua Odoo. Backend sẽ nhận yêu cầu đăng nhập từ webapp, ủy quyền đăng nhập sang Odoo theo URL cấu hình trong biến môi trường, đồng bộ thông tin người dùng về DB Mattermost (tạo mới nếu chưa tồn tại, cập nhật nếu có), sau đó hoàn tất phiên đăng nhập trong Mattermost.

### Mục tiêu
- **Đăng nhập qua Odoo**: Người dùng nhập `username/password` (hoặc email) trên giao diện Mattermost, backend gọi Odoo để xác thực.
- **Đồng bộ user**: Nếu chưa có user tương ứng trong DB Mattermost thì tạo mới; nếu đã có thì cập nhật thông tin cần thiết và đăng nhập.
- **Tối thiểu thay đổi UI**: Tận dụng form login sẵn có; chỉ thay đổi endpoint được gọi.
- **Bảo mật**: Không để lộ thông tin Odoo URL/secret ra client; kiểm soát timeout, retry, và logging an toàn.
- **Có thể rollback**: Cho phép bật/tắt tính năng qua feature flag.

---

## Cấu hình môi trường (backend)

Thêm các biến môi trường sau vào server:

```bash
# Bắt buộc
MM_ODOO_SSO_ENABLED=true                     # Bật/tắt tích hợp Odoo (feature flag)
MM_ODOO_BASE_URL=https://odoo.example.com    # Base URL Odoo được phép gọi

# Lựa chọn phương thức xác thực với Odoo
MM_ODOO_AUTH_METHOD=web_session              # web_session | jsonrpc
MM_ODOO_WEB_AUTH_PATH=/web/session/authenticate
MM_ODOO_JSONRPC_PATH=/jsonrpc

# Tham số Odoo DB (bắt buộc cho jsonrpc, cần cho web_session nếu multi-db)
MM_ODOO_DB=odoo_db_name

# Tuỳ chọn/khuyến nghị
MM_ODOO_TIMEOUT_MS=8000                      # Timeout khi gọi Odoo
MM_ODOO_RETRY=1                              # Số lần retry nhẹ khi lỗi tạm thời
MM_ODOO_CLIENT_ID=mm-backend                 # Nếu Odoo yêu cầu client credentials
MM_ODOO_CLIENT_SECRET=xxxxx                  # Secret tương ứng (giữ bí mật)
MM_ODOO_TLS_INSECURE_SKIP_VERIFY=false       # Chỉ dùng môi trường dev
```

Lưu ý: Tên biến có thể điều chỉnh theo convention của Mattermost server. Nếu hệ thống đã có subsystem cho OAuth/SAML, đây là hướng tích hợp proxy đơn giản không thay đổi core auth flow phức tạp.

---

## Hợp đồng API giữa Webapp ⇄ Backend

Sử dụng endpoint login hiện có của Mattermost nhưng thêm mode Odoo, hoặc tạo endpoint mới dành riêng (khuyến nghị tạo endpoint mới để tránh ảnh hưởng luồng hiện hữu, sau đó có thể alias):

- **POST** `/api/v4/odoo/login`

Request body (JSON):
```json
{
  "identifier": "user@example.com",  
  "password": "plain-text-or-token"
}
```

Response (200):
```json
{
  "user_id": "xxxxxxxxxxxxxxxxxxxx",
  "username": "user",
  "email": "user@example.com",
  "token": "<session-or-personal-access-token>",
  "create": false,
  "updated_fields": ["email", "first_name"]
}
```

Lỗi phổ biến:
- 400: Thiếu tham số
- 401: Sai thông tin đăng nhập từ Odoo
- 409: Xung đột mapping (ví dụ email/username đã thuộc user khác)
- 502/504: Odoo không phản hồi hoặc timeout

---

## Hợp đồng API giữa Backend ⇄ Odoo

Có hai lựa chọn tích hợp phổ biến, tuỳ cấu hình hệ thống:

1) Web session HTTP: `POST ${MM_ODOO_BASE_URL}${MM_ODOO_WEB_AUTH_PATH}` với `Content-Type: application/json`

Request:
```json
{
  "db": "${MM_ODOO_DB}",
  "login": "user@example.com",
  "password": "plain-text"
}
```

Response (200):
```json
{
  "jsonrpc": "2.0",
  "id": null,
  "result": {
    "uid": 123,                      
    "user_context": {"lang": "en_US", "tz": "UTC"},
    "company_id": 1,
    "partner_display_name": "User Example",
    "username": "user@example.com"
  }
}
```

- Cookie: Odoo sẽ thiết lập cookie `session_id` trong header `Set-Cookie`. Nếu cần duy trì phiên cho các call web tiếp theo, backend phải lưu và gửi lại cookie này.
- Thất bại xác thực trả về `result: false` hoặc mã lỗi HTTP kèm thông tin trong trường `error` (tuỳ phiên bản Odoo).

2) JSON-RPC: `POST ${MM_ODOO_BASE_URL}${MM_ODOO_JSONRPC_PATH}` với payload chuẩn JSON-RPC gọi service `common`/method `authenticate`

Request:
```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": {
    "service": "common",
    "method": "authenticate",
    "args": ["${MM_ODOO_DB}", "user@example.com", "plain-text", {}]
  },
  "id": 1
}
```

Response (200):
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": 123     
}
```

- `result` là `uid` (số) nếu thành công; `null` nếu thất bại xác thực.
- Sau khi có `uid`, để lấy thông tin người dùng, gọi JSON-RPC `object`/`execute_kw` trên `res.users` hoặc `res.partner`.

Ví dụ lấy thông tin người dùng (JSON-RPC):
```json
{
  "jsonrpc": "2.0",
  "method": "call",
  "params": {
    "service": "object",
    "method": "execute_kw",
    "args": [
      "${MM_ODOO_DB}",
      123,
      "plain-text",
      "res.users",
      "search_read",
      [[["id", "=", 123]]],
      {"fields": ["id", "name", "login", "email", "partner_id"]}
    ]
  },
  "id": 2
}
```

Gợi ý thực tế:
- Nếu chỉ cần xác thực và không cần giữ phiên Odoo, dùng JSON-RPC `authenticate` là gọn.
- Nếu cần theo sau những call web session (ví dụ `/web/session/get_session_info`), chọn web session để tái sử dụng `session_id`.
- Luôn dùng HTTPS; tránh log mật khẩu; nếu có cơ chế client credential riêng của gateway Odoo, thêm header `Authorization: Bearer ...` vào cả hai luồng.

---

## Thiết kế luồng xử lý

1) Webapp gửi `POST /api/v4/odoo/login` với `identifier` và `password`.
2) Backend kiểm tra `MM_ODOO_SSO_ENABLED` (feature flag). Nếu false, trả 404 hoặc 400 hướng dẫn dùng login thường.
3) Backend gọi Odoo Auth API với timeout/retry đã cấu hình (theo `MM_ODOO_AUTH_METHOD`).
4) Nếu Odoo xác thực thành công:
   - Với web session: đọc `uid` từ `result.uid` và có thể lưu `session_id` nếu cần cho các call theo sau.
   - Với JSON-RPC: `result` là `uid`.
   - Lấy thông tin user bổ sung (name, email) qua JSON-RPC `res.users`/`res.partner` hoặc qua `/web/session/get_session_info` nếu dùng web session.
   - Chuẩn hoá dữ liệu và tìm/ tạo user trong Mattermost như mô tả.
5) Nếu Odoo trả lỗi xác thực: trả 401.
6) Nếu Odoo timeout hoặc lỗi server: trả 502/504, không fallback trừ khi được cấu hình.

Ghi chú mapping:
- Lưu `odoo_user_id` = `uid`.
- Trường email có thể ở `res.users.email` hoặc `res.partner.email` (qua `partner_id`). Nếu email không có, cân nhắc dùng `login`.

---

## Thay đổi backend (dự kiến)

- Cấu hình:
  - Thêm đọc biến môi trường nêu trên vào config server.
  - Thêm validate config khi khởi động.
- HTTP client:
  - Tạo client với timeout, retry và TLS options theo config.
- Endpoint mới:
  - `POST /api/v4/odoo/login` trong module auth hoặc plugin-like package.
  - Validate input (identifier, password).
  - Gọi Odoo, parse response, handle lỗi rõ ràng.
- User service:
  - Hàm `findOrCreateUserFromOdoo(odooProfile)` thực hiện: tìm theo mapping → email → username; tạo mới nếu không có; cập nhật field an toàn.
  - Lưu mapping `odoo_user_id` vào user props hoặc bảng chuyên dụng.
- Session/token:
  - Tái sử dụng cơ chế phát token/thiết lập session có sẵn (như login nội bộ).
- Logging/metrics:
  - Log idempotent (không log password). Thêm metric tỉ lệ thành công, độ trễ, lỗi theo mã.

---

## Thay đổi frontend (dự kiến)

- Sử dụng form login hiện có, nhưng khi feature flag bật trên server, webapp sẽ gọi `/api/v4/odoo/login` thay cho `/api/v4/users/login`.
- Xử lý lỗi 401 (sai thông tin) và 502/504 (dịch vụ Odoo không khả dụng) với thông báo UX thân thiện.
- Không hiển thị thông tin cấu hình Odoo trên client.

Tùy chọn: thêm toggle ẩn (server-driven config qua `/config`) để UI tự động gọi endpoint Odoo khi bật.

---

## Bảo mật

- Chỉ gọi Odoo qua HTTPS. Không bao giờ log `password` hoặc `client_secret`.
- Rate limiting và lockout theo chính sách hiện có của Mattermost, áp dụng ở lớp trước khi gọi Odoo để tránh brute-force.
- Ràng buộc CORS vẫn giữ nguyên, endpoint chỉ dành cho webapp hợp lệ.
- TTL session/token giữ nguyên chuẩn của hệ thống.

---

## Xử lý cạnh

- Email thay đổi bên Odoo: cập nhật nếu không gây xung đột; nếu xung đột, trả 409 và hướng dẫn quy trình hợp nhất.
- User bị deactivated trên Odoo: chặn đăng nhập, có thể auto-deactivate trong MM nếu policy cho phép.
- Odoo tạm thời lỗi: hiển thị thông báo và không fallback trừ khi có `MM_ODOO_FALLBACK_LOCAL=true` (tuỳ chọn) và user cũng có mật khẩu local.

---

## Kiểm thử

- Unit tests:
  - Parser/validator response Odoo
  - `findOrCreateUserFromOdoo` với các ca: mới, đã tồn tại theo mapping/email/username, xung đột email
- Integration tests (backend):
  - Mock Odoo trả 200/401/5xx/timeout; kiểm tra mã lỗi và side-effects DB
- E2E (cypress):
  - Đăng nhập thành công; tạo mới user
  - Sai mật khẩu → 401
  - Odoo timeout → thông báo lỗi

---

## Rollout & vận hành

- Giai đoạn 1: ẩn sau feature flag `MM_ODOO_SSO_ENABLED=false` (default)
- Giai đoạn 2: bật ở môi trường staging với Odoo giả lập/mocking
- Giai đoạn 3: canary một nhóm nhỏ người dùng
- Giám sát metrics, error logs; chuẩn bị rollback bằng cách tắt flag

---

## Pseudo-code backend (minh hoạ)

```go
func OdooLoginHandler(w http.ResponseWriter, r *http.Request) {
    if !cfg.OdooEnabled { return notFound() }
    req := parseLoginRequest(r)
    if err := validate(req); err != nil { return badRequest(err) }

    odooResp, err := odooClient.Authenticate(req.Identifier, req.Password)
    if err != nil { return upstreamError(err) }
    if !odooResp.IsActive { return unauthorized("inactive user") }

    user, created, updatedFields, err := userSvc.FindOrCreateFromOdoo(odooResp)
    if err != nil { return conflictOrServerError(err) }

    token, err := sessionSvc.IssueForUser(user.Id)
    if err != nil { return serverError(err) }

    respondJSON(w, 200, map[string]any{
        "user_id": user.Id,
        "username": user.Username,
        "email": user.Email,
        "token": token,
        "create": created,
        "updated_fields": updatedFields,
    })
}
```

---

## Công việc thực thi (tóm tắt)

- Backend: config, HTTP client Odoo, endpoint `/api/v4/odoo/login`, mapping user, session, logs/metrics.
- Frontend: chuyển endpoint login khi flag bật, xử lý thông báo lỗi.
- Test: unit/integration/E2E; tài liệu vận hành và rollback.

```text
Trạng thái: Bản kế hoạch v1 – sẵn sàng triển khai.
```
