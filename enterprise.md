## Enterprise License Audit & Plan

### Mục tiêu
- Tổng hợp các điểm kiểm tra license ở frontend/backend.
- Chuẩn hoá logic phân quyền theo license (Free/Cloud/Starter/Professional/Enterprise/Enterprise Advanced/Entry).
- Đề xuất kế hoạch chuẩn hoá, test, và tài liệu hoá.

### Nguồn dữ liệu license
- Backend trả về client license qua endpoint: `GET /api/v4/license/client?format=old` (file `server/channels/api4/license.go`).
- Server giữ `License` tại: `platform.PlatformService.licenseValue` với các API:
  - `Server.License()` → `PlatformService.License()` (file `server/channels/app/license.go`, `server/channels/app/platform/license.go`).
  - Load/Set/Validate/Save license: `LoadLicense/SetLicense/ValidateAndSetLicenseBytes/SaveLicense`.
- Mô hình dữ liệu server: `server/public/model/license.go` (`License`, `Features`, các helper như `IsCloud`, `IsTrialLicense`, `MinimumEnterpriseLicense`, ...).

### Frontend: nơi đọc/kiểm tra license
- Selector trung tâm: `getLicense(state)` trong `webapp/channels/src/packages/mattermost-redux/src/selectors/entities/general.ts`.
- Các tiện ích/logic:
  - `webapp/channels/src/utils/license_utils.ts`:
    - `isCloudLicense`, `isEnterpriseLicense`, `isEnterpriseOrCloudOrSKUStarterFree`, `isMinimumProfessionalLicense`, `isMinimumEnterpriseLicense`, `isMinimumEnterpriseAdvancedLicense`.
    - Xử lý hết hạn/đang hết hạn: `isLicenseExpiring`, `isLicenseExpired`, `isLicensePastGracePeriod`, `daysToLicenseExpire` (bỏ qua Cloud).
- Nhiều component/selector tham chiếu trực tiếp `getLicense(state)` để gate tính năng (MFA, LDAP groups, compliance, banner, menu, admin console, v.v.).

### Backend: nơi kiểm tra/thi hành license
- Router/API: `server/channels/api4/license.go`
  - `POST /license` thêm license, kiểm tra trial eligibility, lưu store, reload config, invalidate cache.
  - `DELETE /license` xoá license.
  - `POST /trial-license` xin trial license (kiểm tra quyền, eligibility).
  - `GET /license/client` trả về client license (đã lọc theo quyền đọc license info).
  - `GET /trial-license/prev` lấy trial trước đó (qua LicenseManager).
- Service layer: `server/channels/app/platform/license.go`
  - Nguồn license: ENV `MM_LICENSE` → DB `SystemActiveLicenseId` → file `LicenseFileLocation`.
  - `SetLicense` đồng bộ `clientLicense` và phát sự kiện tới listeners.
  - `SaveLicense` kiểm tra số user, hết hạn, dừng/khởi động workers/schedulers, lưu `SystemActiveLicenseId` và record.
  - `RequestTrialLicense` gọi CWS, cài license, reload/invalidate cache.
- App layer: `server/channels/app/license.go`
  - Phơi bày `Server.License()`, `Load/Save/Set/Validate/ClientLicense()` và logic request trial có mở rộng trường.
- Model: `server/public/model/license.go`
  - Xác định tiers: Professional(10), Enterprise(20), Enterprise Advanced(30), Entry(30).
  - Helper: `IsCloud`, `IsTrialLicense`, `IsExpired`, `IsPastGracePeriod`, `Minimum*License`.

### Phân loại license và gating chính
- Cloud: `license.Features.Cloud == true` → FE bỏ qua cảnh báo hết hạn; BE có luồng invalidate cache on change.
- Self-hosted SKU:
  - Starter (Free): FE coi "Starter Free" qua `isEnterpriseReady && license.IsLicensed === 'false'` hoặc `SelfHostedProducts === STARTER`.
  - Professional ≥10: yêu cầu tối thiểu để bật một số features (ví dụ Shared Channels/Remote Cluster cũng mở nếu tối thiểu Pro).
  - Enterprise ≥20: coi là Enterprise; nhiều tính năng admin nâng cao yêu cầu Enterprise.
  - Enterprise Advanced ≥30.
  - Entry: giấy phép đặc biệt do `FeatureFlags.EnableMattermostEntry` có thể mở enterprise features không có license DB.

### Điểm cần chuẩn hoá/đồng bộ FE-BE
- Định nghĩa tier: FE dùng `getLicenseTier` theo `LicenseSkus` và BE dùng `LicenseToLicenseTier`; cần mapping một-một.
- Kiểm tra "Starter Free" ở FE: đồng bộ với backend semantics (SelfHostedProducts vs IsLicensed === 'false').
- Trạng thái Cloud vs Self-hosted: FE dùng `license.Cloud === 'true'`; BE dùng `License.IsCloud()` theo `Features.Cloud`.
- Trải nghiệm trial: FE hiển thị banner/modals; BE enforce `CanStartTrial` và sanctioned trial.

### Kế hoạch thực hiện
1) Kiểm kê và nhóm hoá gating FE
   - Tạo danh sách feature → điều kiện license/sku/tier/flag.
   - Chuẩn hoá dùng các helper trong `license_utils.ts` thay vì check rải rác `license.IsLicensed === 'true'`.

2) Chuẩn hoá tầng selector
   - Tạo các selector có tên rõ nghĩa: `selectIsCloud`, `selectIsEnterpriseTier`, `selectIsProfessionalOrHigher`, `selectIsStarterFree` (bọc helpers).
   - Dần thay thế các nơi đọc trực tiếp `getLicense(state)` để logic tập trung.

3) Backend validation pass
   - Rà `SaveLicense`, `LoadLicense`, `RequestTrialLicense` để đảm bảo error path rõ, log đủ, và cache invalidation nhất quán.
   - Xác thực mapping SKU ↔ tier đồng nhất với FE.

4) Test plan
   - Unit FE: `license_utils.test.ts` bổ sung matrix test cho tất cả SKU/tier (Starter/Pro/Enterprise/Advanced/Entry) và Cloud/Trial.
   - Unit BE: test `Minimum*License`, `IsCloud/IsTrial/IsExpired`, `RequestTrialLicense` path (sanctioned trial, eligibility fail/pass).
   - E2E: cài/ghi/xoá license, verify gating trong UI chính (Admin Console sections, LDAP/MFA, Marketplace enterprise plugins, Shared Channels).

5) Tài liệu
   - Mục "Licensing & Feature Gating" trong docs nội bộ: giải thích trường trong client license, SKU/tier, trình tự load license, và các selector FE nên dùng.

### Checklist thực thi
- [ ] Thêm/chuẩn hoá selector license tier ở FE.
- [ ] Refactor các component chính sang dùng selector thay vì check ad-hoc.
- [ ] Bổ sung test matrix FE/BE cho SKU và Cloud/Trial/Expiry.
- [ ] Đối chiếu mapping SKU giữa FE `utils/constants` và BE `model.LicenseToLicenseTier`.
- [ ] Viết docs nội bộ, link đến file này.


