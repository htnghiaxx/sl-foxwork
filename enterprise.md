# Kế hoạch loại bỏ hạn chế Enterprise về Clustering

## Tổng quan

Tài liệu này mô tả kế hoạch chi tiết để loại bỏ các hạn chế Enterprise liên quan đến clustering trong Mattermost, cho phép tất cả người dùng sử dụng tính năng clustering mà không cần license Enterprise.

## Các tính năng clustering hiện tại bị hạn chế Enterprise

### 1. Redis Clustering
- **Vị trí**: `server/platform/services/cache/provider.go`
- **Hạn chế**: Yêu cầu license với `Features.Cluster = true`
- **Code hiện tại**:
  ```go
  if (license == nil || !*license.Features.Cluster) && *cacheConfig.CacheType == model.CacheTypeRedis && !ps.forceEnableRedis {
      return nil, fmt.Errorf("Redis cannot be used in an instance without a license or a license without clustering")
  }
  ```

### 2. WebSocket Clustering
- **Vị trí**: `server/channels/app/platform/cluster.go`
- **Hạn chế**: Yêu cầu license với `Features.Cluster = true`
- **Code hiện tại**:
  ```go
  if ps.License() != nil && *ps.Config().ClusterSettings.Enable && ps.clusterIFace != nil {
      return ps.clusterIFace.IsLeader()
  }
  ```

### 3. Cluster Settings
- **Vị trí**: `server/public/model/config.go`
- **Hạn chế**: `ClusterSettings.Enable` mặc định là `false`
- **Access control**: `environment_high_availability,write_restrictable,cloud_restrictable`

## Kế hoạch thực hiện

### Phase 1: Chuẩn bị và phân tích (Tuần 1-2)

#### 1.1 Audit toàn bộ codebase
- [ ] Tìm tất cả các kiểm tra license liên quan đến clustering
- [ ] Xác định các dependencies và side effects
- [ ] Đánh giá tác động đến performance và security

#### 1.2 Tạo feature flags
- [ ] Thêm `EnableOpenSourceClustering` trong config
- [ ] Tạo migration path cho existing deployments
- [ ] Thiết lập monitoring và logging

### Phase 2: Loại bỏ hạn chế Redis (Tuần 3-4)

#### 2.1 Sửa đổi Redis provider
```go
// File: server/platform/services/cache/provider.go
// Thay đổi từ:
if (license == nil || !*license.Features.Cluster) && *cacheConfig.CacheType == model.CacheTypeRedis && !ps.forceEnableRedis {
    return nil, fmt.Errorf("Redis cannot be used in an instance without a license or a license without clustering")
}

// Thành:
if *cacheConfig.CacheType == model.CacheTypeRedis {
    // Redis luôn available, không cần kiểm tra license
    ps.cacheProvider, err = cache.NewRedisProvider(...)
}
```

#### 2.2 Cập nhật config validation
```go
// File: server/public/model/config.go
// Thay đổi ClusterSettings access control
type ClusterSettings struct {
    Enable *bool `access:"write_restrictable"` // Loại bỏ environment_high_availability
    // ... other fields
}
```

### Phase 3: Loại bỏ hạn chế WebSocket Clustering (Tuần 5-6)

#### 3.1 Sửa đổi cluster logic
```go
// File: server/channels/app/platform/cluster.go
func (ps *PlatformService) IsLeader() bool {
    // Loại bỏ kiểm tra license
    if *ps.Config().ClusterSettings.Enable && ps.clusterIFace != nil {
        return ps.clusterIFace.IsLeader()
    }
    return true
}
```

#### 3.2 Cập nhật WebSocket publishing
```go
// File: server/channels/app/platform/cluster.go
func (ps *PlatformService) Publish(message *model.WebSocketEvent) {
    ps.PublishSkipClusterSend(message)
    
    // Luôn gửi cluster message nếu clustering enabled
    if ps.clusterIFace != nil && *ps.Config().ClusterSettings.Enable {
        // ... existing cluster logic
    }
}
```

### Phase 4: Cập nhật Cluster Settings (Tuần 7-8)

#### 4.1 Thay đổi default values
```go
// File: server/public/model/config.go
func (s *ClusterSettings) SetDefaults() {
    if s.Enable == nil {
        s.Enable = NewPointer(true) // Thay đổi từ false thành true
    }
    // ... other defaults
}
```

#### 4.2 Loại bỏ access restrictions
```go
// Loại bỏ environment_high_availability từ tất cả cluster settings
type ClusterSettings struct {
    Enable                      *bool   `access:"write_restrictable"`
    ClusterName                 *string `access:"write_restrictable"`
    NetworkInterface           *string `access:"write_restrictable"`
    // ... other fields
}
```

### Phase 5: Cập nhật UI và Documentation (Tuần 9-10)

#### 5.1 Cập nhật Admin Console
- [ ] Loại bỏ warnings về Enterprise license cho clustering
- [ ] Cập nhật help text và tooltips
- [ ] Thêm migration guide cho existing users

#### 5.2 Cập nhật Documentation
- [ ] Cập nhật installation guides
- [ ] Thêm clustering best practices
- [ ] Cập nhật API documentation

### Phase 6: Testing và Validation (Tuần 11-12)

#### 6.1 Unit Tests
- [ ] Cập nhật tất cả tests liên quan đến clustering
- [ ] Thêm tests cho open source clustering
- [ ] Test migration scenarios

#### 6.2 Integration Tests
- [ ] Test multi-instance deployments
- [ ] Test Redis clustering without license
- [ ] Test WebSocket clustering without license
- [ ] Performance testing

#### 6.3 End-to-End Tests
- [ ] Test complete clustering setup
- [ ] Test failover scenarios
- [ ] Test load balancing

## Migration Strategy

### Cho Existing Deployments

#### 1. Automatic Migration
```go
// Thêm migration logic
func (s *ClusterSettings) MigrateFromEnterprise() {
    // Nếu đang sử dụng clustering với license, giữ nguyên
    // Nếu không có license nhưng có cluster config, enable clustering
    if s.Enable == nil && hasClusterConfig() {
        s.Enable = NewPointer(true)
    }
}
```

#### 2. Configuration Updates
- [ ] Tự động enable clustering nếu detect Redis config
- [ ] Thêm warnings cho users về changes
- [ ] Provide rollback mechanism

### Cho New Deployments
- [ ] Clustering enabled by default
- [ ] Simplified setup process
- [ ] Clear documentation về benefits

## Risk Assessment và Mitigation

### High Risk
1. **Performance Impact**
   - **Risk**: Clustering có thể impact performance
   - **Mitigation**: Comprehensive performance testing, monitoring

2. **Security Concerns**
   - **Risk**: Clustering có thể tạo security vulnerabilities
   - **Mitigation**: Security audit, encryption requirements

### Medium Risk
1. **Configuration Complexity**
   - **Risk**: Users có thể misconfigure clustering
   - **Mitigation**: Better defaults, validation, documentation

2. **Backward Compatibility**
   - **Risk**: Breaking changes cho existing deployments
   - **Mitigation**: Gradual migration, feature flags

### Low Risk
1. **Documentation Updates**
   - **Risk**: Outdated documentation
   - **Mitigation**: Comprehensive documentation review

## Success Metrics

### Technical Metrics
- [ ] 100% clustering features available without license
- [ ] Zero performance regression
- [ ] All existing tests pass
- [ ] New clustering tests pass

### User Metrics
- [ ] Increased adoption of clustering features
- [ ] Reduced support tickets về licensing
- [ ] Positive feedback từ community

## Timeline

| Phase | Duration | Deliverables |
|-------|----------|--------------|
| Phase 1 | 2 tuần | Audit report, feature flags |
| Phase 2 | 2 tuần | Redis clustering without license |
| Phase 3 | 2 tuần | WebSocket clustering without license |
| Phase 4 | 2 tuần | Updated cluster settings |
| Phase 5 | 2 tuần | Updated UI and docs |
| Phase 6 | 2 tuần | Testing and validation |

**Total Duration**: 12 tuần (3 tháng)

## Rollback Plan

Nếu có issues nghiêm trọng:

1. **Immediate Rollback**
   - Revert code changes
   - Restore license checks
   - Notify users

2. **Gradual Rollback**
   - Disable feature flags
   - Provide migration path back
   - Monitor impact

3. **Communication**
   - Notify community
   - Update documentation
   - Provide support

## Conclusion

Kế hoạch này sẽ loại bỏ hoàn toàn hạn chế Enterprise về clustering, cho phép tất cả users sử dụng tính năng clustering mà không cần license. Điều này sẽ:

- Tăng adoption của Mattermost
- Cải thiện user experience
- Giảm complexity trong deployment
- Tăng community engagement

Việc thực hiện sẽ được thực hiện một cách cẩn thận với comprehensive testing và monitoring để đảm bảo stability và performance.
