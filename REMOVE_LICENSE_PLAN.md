# Kế hoạch loại bỏ hoàn toàn hệ thống License và Cloud Connections

## Tổng quan

Tài liệu này mô tả kế hoạch chi tiết để loại bỏ hoàn toàn hệ thống license và các kết nối cloud trong Mattermost, biến nó thành một sản phẩm hoàn toàn open source không có hạn chế.

## Phân tích hiện tại

### 1. Các kết nối Cloud cần loại bỏ

#### 1.1 Cloud Settings
```go
// server/public/model/config.go
type CloudSettings struct {
    CWSURL                *string `access:"write_restrictable"`
    CWSAPIURL             *string `access:"write_restrictable"`
    CWSMock               *bool   `access:"write_restrictable"`
    Disable               *bool   `access:"write_restrictable,cloud_restrictable"`
    PreviewModalBucketURL *string `access:"write_restrictable"`
}

// Default URLs cần loại bỏ:
CloudSettingsDefaultCwsURL        = "https://customers.mattermost.com"
CloudSettingsDefaultCwsAPIURL     = "https://portal.internal.prod.cloud.mattermost.com"
CloudSettingsDefaultCwsURLTest    = "https://portal.test.cloud.mattermost.com"
CloudSettingsDefaultCwsAPIURLTest = "https://api.internal.test.cloud.mattermost.com"
```

#### 1.2 License Server Connections
```go
// server/channels/app/platform/license.go
func (ps *PlatformService) RequestTrialLicense(trialRequest *model.TrialLicenseRequest) *model.AppError {
    resp, err := http.Post(ps.getRequestTrialURL(), "application/json", bytes.NewBuffer(trialRequestJSON))
    // Kết nối đến: https://customers.mattermost.com/api/v1/trials
}
```

### 2. Hệ thống License cần loại bỏ

#### 2.1 License Validation
- `utils.LicenseValidator.ValidateLicense()`
- `utils.LicenseValidator.LicenseFromBytes()`
- JWT validation cho license renewal

#### 2.2 License Storage
- Database storage cho license records
- File system storage cho license files
- Environment variable loading (`MM_LICENSE`)

#### 2.3 License Features
- Tất cả `Features.*` checks trong code
- Enterprise feature restrictions
- User count limitations

## Kế hoạch thực hiện

### Phase 1: Tạo Open Source License System (Tuần 1-2)

#### 1.1 Tạo Open Source License
```go
// server/public/model/opensource_license.go
type OpenSourceLicense struct {
    Id        string    `json:"id"`
    IssuedAt  int64     `json:"issued_at"`
    ExpiresAt int64     `json:"expires_at"`
    Features  *Features `json:"features"`
}

func NewOpenSourceLicense() *OpenSourceLicense {
    return &OpenSourceLicense{
        Id:        "opensource-permanent",
        IssuedAt:  model.GetMillis(),
        ExpiresAt: model.GetMillis() + (100 * 365 * 24 * 60 * 60 * 1000), // 100 years
        Features: &Features{
            Users:                     NewPointer(999999999), // Unlimited
            LDAP:                      NewPointer(true),
            LDAPGroups:                NewPointer(true),
            MFA:                       NewPointer(true),
            GoogleOAuth:               NewPointer(true),
            Office365OAuth:            NewPointer(true),
            OpenId:                    NewPointer(true),
            Compliance:                NewPointer(true),
            Cluster:                   NewPointer(true), // Enable clustering
            Metrics:                   NewPointer(true),
            MHPNS:                     NewPointer(true),
            SAML:                      NewPointer(true),
            Elasticsearch:             NewPointer(true), // Enable search
            Announcement:              NewPointer(true),
            ThemeManagement:           NewPointer(true),
            EmailNotificationContents: NewPointer(true),
            DataRetention:             NewPointer(true),
            MessageExport:             NewPointer(true),
            CustomPermissionsSchemes:  NewPointer(true),
            CustomTermsOfService:      NewPointer(true),
            GuestAccounts:             NewPointer(true),
            GuestAccountsPermissions:  NewPointer(true),
            IDLoadedPushNotifications: NewPointer(true),
            LockTeammateNameDisplay:   NewPointer(true),
            EnterprisePlugins:         NewPointer(true),
            AdvancedLogging:           NewPointer(true),
            Cloud:                     NewPointer(false), // Disable cloud
            SharedChannels:            NewPointer(true),
            RemoteClusterService:      NewPointer(true),
            OutgoingOAuthConnections:  NewPointer(true),
            FutureFeatures:            NewPointer(true),
        },
    }
}
```

#### 1.2 Tạo Open Source License Manager
```go
// server/einterfaces/opensource_license.go
type OpenSourceLicenseInterface interface {
    GetLicense() *model.License
    IsLicensed() bool
    HasFeature(feature string) bool
}

type OpenSourceLicenseManager struct {
    license *model.License
}

func NewOpenSourceLicenseManager() *OpenSourceLicenseManager {
    return &OpenSourceLicenseManager{
        license: NewOpenSourceLicense(),
    }
}

func (osm *OpenSourceLicenseManager) GetLicense() *model.License {
    return osm.license
}

func (osm *OpenSourceLicenseManager) IsLicensed() bool {
    return true // Always licensed in open source
}

func (osm *OpenSourceLicenseManager) HasFeature(feature string) bool {
    return true // All features enabled in open source
}
```

### Phase 2: Loại bỏ Cloud Connections (Tuần 3-4)

#### 2.1 Loại bỏ Cloud Settings
```go
// server/public/model/config.go
// Xóa hoàn toàn CloudSettings struct
// Xóa tất cả references đến CloudSettings
// Xóa CloudSettingsDefaultCwsURL, CloudSettingsDefaultCwsAPIURL, etc.
```

#### 2.2 Loại bỏ Trial License Request
```go
// server/channels/app/platform/license.go
// Xóa hoàn toàn RequestTrialLicense function
// Xóa getRequestTrialURL function
// Xóa TrialLicenseRequest struct
```

#### 2.3 Loại bỏ Cloud API endpoints
```go
// server/channels/api4/cloud.go
// Xóa hoàn toàn file này
// Xóa tất cả cloud-related API endpoints
```

### Phase 3: Thay thế License System (Tuần 5-6)

#### 3.1 Sửa đổi license.go
```go
// server/channels/app/platform/license.go
// Thay thế toàn bộ nội dung:

package platform

import (
    "github.com/mattermost/mattermost/server/public/model"
    "github.com/mattermost/mattermost/server/v8/einterfaces"
)

func (ps *PlatformService) LicenseManager() einterfaces.LicenseInterface {
    return ps.licenseManager
}

func (ps *PlatformService) SetLicenseManager(impl einterfaces.LicenseInterface) {
    ps.licenseManager = impl
}

func (ps *PlatformService) License() *model.License {
    return ps.licenseManager.GetLicense()
}

func (ps *PlatformService) LoadLicense() {
    // Không cần load gì, license luôn available
    ps.logger.Info("Open source license loaded, all features enabled.")
}

func (ps *PlatformService) SetLicense(license *model.License) bool {
    // Không cần set license, luôn có open source license
    return true
}

func (ps *PlatformService) ClientLicense() map[string]string {
    return map[string]string{
        "IsLicensed": "true",
        "Users": "999999999",
        "LDAP": "true",
        "LDAPGroups": "true",
        "MFA": "true",
        "GoogleOAuth": "true",
        "Office365OAuth": "true",
        "OpenId": "true",
        "Compliance": "true",
        "Cluster": "true",
        "Metrics": "true",
        "MHPNS": "true",
        "SAML": "true",
        "Elasticsearch": "true",
        "Announcement": "true",
        "ThemeManagement": "true",
        "EmailNotificationContents": "true",
        "DataRetention": "true",
        "MessageExport": "true",
        "CustomPermissionsSchemes": "true",
        "CustomTermsOfService": "true",
        "GuestAccounts": "true",
        "GuestAccountsPermissions": "true",
        "IDLoadedPushNotifications": "true",
        "LockTeammateNameDisplay": "true",
        "EnterprisePlugins": "true",
        "AdvancedLogging": "true",
        "Cloud": "false",
        "SharedChannels": "true",
        "RemoteClusterService": "true",
        "OutgoingOAuthConnections": "true",
        "FutureFeatures": "true",
    }
}

func (ps *PlatformService) GetSanitizedClientLicense() map[string]string {
    return ps.ClientLicense()
}

// Xóa tất cả functions khác liên quan đến license validation, storage, etc.
```

#### 3.2 Cập nhật Platform Service
```go
// server/channels/app/platform/service.go
// Trong NewPlatformService, thay thế:
func NewPlatformService(...) (*PlatformService, error) {
    // ... existing code ...
    
    // Thay vì load license, set open source license manager
    ps.licenseManager = NewOpenSourceLicenseManager()
    
    // ... rest of initialization ...
}
```

### Phase 4: Loại bỏ License Dependencies (Tuần 7-8)

#### 4.1 Loại bỏ License Validator
```go
// server/channels/utils/license.go
// Xóa hoàn toàn file này
// Xóa tất cả license validation logic
```

#### 4.2 Loại bỏ License Storage
```go
// server/channels/store/store.go
// Xóa License() store interface
// Xóa tất cả license storage methods
```

#### 4.3 Loại bỏ License API
```go
// server/channels/api4/license.go
// Xóa hoàn toàn file này
// Xóa tất cả license-related API endpoints
```

### Phase 5: Cập nhật Feature Checks (Tuần 9-10)

#### 5.1 Thay thế tất cả License Checks
```go
// Thay thế tất cả patterns như:
if license := ps.License(); license == nil || !*license.Features.Cluster {
    return nil
}

// Thành:
// Không cần check gì, tất cả features đều enabled
```

#### 5.2 Cập nhật Config Validation
```go
// server/public/model/config.go
// Loại bỏ tất cả license-related validation
// Enable tất cả features by default
```

### Phase 6: Cập nhật UI và Documentation (Tuần 11-12)

#### 6.1 Loại bỏ License UI
- Xóa tất cả license-related UI components
- Xóa license settings pages
- Xóa trial license requests

#### 6.2 Cập nhật Documentation
- Cập nhật installation guides
- Loại bỏ references đến license requirements
- Thêm open source feature documentation

## Code Changes Chi tiết

### 1. Files cần xóa hoàn toàn:
- `server/channels/api4/cloud.go`
- `server/channels/api4/license.go`
- `server/channels/utils/license.go`
- `server/channels/store/license.go` (nếu có)

### 2. Files cần sửa đổi lớn:
- `server/channels/app/platform/license.go` - Thay thế hoàn toàn
- `server/public/model/config.go` - Loại bỏ CloudSettings
- `server/public/model/license.go` - Thêm OpenSourceLicense
- `server/einterfaces/license.go` - Thêm OpenSourceLicenseInterface

### 3. Files cần sửa đổi nhỏ:
- Tất cả files có license checks cần được cập nhật
- UI components liên quan đến license
- Documentation files

## Migration Strategy

### 1. Backward Compatibility
- Giữ nguyên API structure
- `ClientLicense()` vẫn trả về map như cũ
- `License()` vẫn trả về *model.License

### 2. Configuration Migration
- Tự động enable tất cả features
- Loại bỏ cloud settings từ config
- Không cần migration database

### 3. User Communication
- Thông báo về việc loại bỏ license requirements
- Hướng dẫn về new open source features
- Migration guide cho existing deployments

## Testing Strategy

### 1. Unit Tests
- Test open source license manager
- Test tất cả features hoạt động without license
- Test config migration

### 2. Integration Tests
- Test clustering without license
- Test search without license
- Test tất cả enterprise features

### 3. End-to-End Tests
- Test complete setup without license
- Test feature availability
- Test performance impact

## Risk Assessment

### High Risk
1. **Breaking Changes**
   - Risk: API changes có thể break existing integrations
   - Mitigation: Maintain API compatibility

2. **Performance Impact**
   - Risk: Loại bỏ license checks có thể impact performance
   - Mitigation: Benchmark testing

### Medium Risk
1. **Feature Availability**
   - Risk: Một số features có thể không hoạt động đúng
   - Mitigation: Comprehensive testing

2. **Configuration Changes**
   - Risk: Config changes có thể break existing setups
   - Mitigation: Automatic migration

### Low Risk
1. **Documentation Updates**
   - Risk: Outdated documentation
   - Mitigation: Comprehensive documentation review

## Success Metrics

### Technical Metrics
- [ ] 100% features available without license
- [ ] Zero cloud connections
- [ ] All existing tests pass
- [ ] No performance regression

### User Metrics
- [ ] Simplified setup process
- [ ] No license-related support tickets
- [ ] Increased adoption

## Timeline

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| Phase 1 | 2 tuần | Open source license system |
| Phase 2 | 2 tuần | Remove cloud connections |
| Phase 3 | 2 tuần | Replace license system |
| Phase 4 | 2 tuần | Remove license dependencies |
| Phase 5 | 2 tuần | Update feature checks |
| Phase 6 | 2 tuần | Update UI and docs |

**Total Duration**: 12 tuần (3 tháng)

## Rollback Plan

Nếu có issues nghiêm trọng:

1. **Immediate Rollback**
   - Revert tất cả changes
   - Restore original license system
   - Notify users

2. **Gradual Rollback**
   - Disable open source features
   - Restore license requirements
   - Provide migration path

3. **Communication**
   - Notify community
   - Update documentation
   - Provide support

## Conclusion

Kế hoạch này sẽ loại bỏ hoàn toàn hệ thống license và cloud connections, biến Mattermost thành một sản phẩm hoàn toàn open source với tất cả features available mà không cần license. Điều này sẽ:

- Loại bỏ hoàn toàn vendor lock-in
- Tăng adoption và community engagement
- Đơn giản hóa deployment và maintenance
- Tạo ra một sản phẩm thực sự open source

Việc thực hiện sẽ được thực hiện một cách cẩn thận với comprehensive testing để đảm bảo stability và backward compatibility.
