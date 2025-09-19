#!/bin/bash

# Script to remove all remaining license checks from Mattermost

echo "Removing license checks from Mattermost..."

# Remove license checks from config/client.go
sed -i '' 's/if \*license\.Features\.Compliance {/\/\/ Open source license always has Compliance enabled\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if \*license\.Features\.SAML {/\/\/ Open source license always has SAML enabled\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if \*license\.Features\.Announcement {/\/\/ Open source license always has Announcement enabled\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if \*license\.Features\.ThemeManagement {/\/\/ Open source license always has ThemeManagement enabled\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if \*license\.Features\.DataRetention {/\/\/ Open source license always has DataRetention enabled\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if license\.HasSharedChannels() {/\/\/ Open source license always has SharedChannels enabled\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if model\.MinimumProfessionalLicense(license) {/\/\/ Open source license always has Professional features\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if model\.MinimumEnterpriseLicense(license) {/\/\/ Open source license always has Enterprise features\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go
sed -i '' 's/if model\.MinimumEnterpriseAdvancedLicense(license) {/\/\/ Open source license always has Enterprise Advanced features\n\t{/g' /Users/nghia/DEV/mattermost/server/config/client.go

# Remove license checks from other files
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/admin.go
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/platform/service.go
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/security_update_check.go
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/email/email.go

# Remove license checks from user.go
sed -i '' 's/license := c\.App\.Channels()\.License(); license == nil || !\*license\.Features\.CustomPermissionsSchemes/\/\/ Open source license always has CustomPermissionsSchemes enabled/g' /Users/nghia/DEV/mattermost/server/channels/api4/user.go

# Remove license checks from pluginapi
sed -i '' 's/license != nil/true/g' /Users/nghia/DEV/mattermost/server/public/pluginapi/license.go
sed -i '' 's/license == nil/false/g' /Users/nghia/DEV/mattermost/server/public/pluginapi/license.go

# Remove license checks from model files
sed -i '' 's/license != nil && LicenseToLicenseTier\[license\.SkuShortName\] >= ProfessionalTier/true/g' /Users/nghia/DEV/mattermost/server/public/model/license.go
sed -i '' 's/license != nil && LicenseToLicenseTier\[license\.SkuShortName\] >= EnterpriseTier/true/g' /Users/nghia/DEV/mattermost/server/public/model/license.go
sed -i '' 's/license != nil && LicenseToLicenseTier\[license\.SkuShortName\] >= EnterpriseAdvancedTier/true/g' /Users/nghia/DEV/mattermost/server/public/model/license.go

# Remove license checks from metrics
sed -i '' 's/(license != nil && \*license\.Features\.Metrics) || (model\.BuildNumber == "dev")/true/g' /Users/nghia/DEV/mattermost/server/enterprise/metrics/metrics.go

# Remove license checks from jobs
sed -i '' 's/license != nil && license\.Features != nil && \*license\.Features\.Cloud/true/g' /Users/nghia/DEV/mattermost/server/channels/jobs/notify_admin/worker.go
sed -i '' 's/license != nil && \*license\.Features\.Cloud/true/g' /Users/nghia/DEV/mattermost/server/channels/jobs/notify_admin/scheduler.go
sed -i '' 's/license != nil && license\.Limits != nil && license\.Limits\.PostHistory > 0/true/g' /Users/nghia/DEV/mattermost/server/channels/jobs/last_accessible_post/worker.go
sed -i '' 's/license != nil && license\.Limits != nil && license\.Limits\.PostHistory > 0/true/g' /Users/nghia/DEV/mattermost/server/channels/jobs/last_accessible_post/scheduler.go
sed -i '' 's/license != nil && \*license\.Features\.Cloud/true/g' /Users/nghia/DEV/mattermost/server/channels/jobs/last_accessible_file/worker.go
sed -i '' 's/license != nil && \*license\.Features\.Cloud/true/g' /Users/nghia/DEV/mattermost/server/channels/jobs/last_accessible_file/scheduler.go
sed -i '' 's/model\.BuildEnterpriseReady == "true" && license == nil/false/g' /Users/nghia/DEV/mattermost/server/channels/jobs/hosted_purchase_screening/scheduler.go

# Remove license checks from app files
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/users/profile_picture.go
sed -i '' 's/license != nil && \*license\.Features\.AdvancedLogging/true/g' /Users/nghia/DEV/mattermost/server/channels/app/server.go
sed -i '' 's/license == nil || license\.Limits == nil || license\.Limits\.PostHistory == 0/false/g' /Users/nghia/DEV/mattermost/server/channels/app/post.go
sed -i '' 's/license == nil || !\*license\.Features\.LDAPGroups/true/g' /Users/nghia/DEV/mattermost/server/channels/app/plugin_api.go
sed -i '' 's/license != nil && license\.HasEnterpriseMarketplacePlugins()/true/g' /Users/nghia/DEV/mattermost/server/channels/app/plugin.go
sed -i '' 's/license != nil && license\.IsCloud()/false/g' /Users/nghia/DEV/mattermost/server/channels/app/plugin.go
sed -i '' 's/license := ps\.License(); license != nil/\/\/ Open source license always available/g' /Users/nghia/DEV/mattermost/server/channels/app/platform/support_packet.go
sed -i '' 's/license != nil && \*license\.Features\.AdvancedLogging/true/g' /Users/nghia/DEV/mattermost/server/channels/app/platform/log.go
sed -i '' 's/license != nil/true/g' /Users/nghia/DEV/mattermost/server/channels/app/notify_admin.go
sed -i '' 's/license := a\.Srv()\.License(); license != nil && \*license\.Features\.EmailNotificationContents/\/\/ Open source license always has EmailNotificationContents enabled/g' /Users/nghia/DEV/mattermost/server/channels/app/notification_email.go
sed -i '' 's/license := a\.Srv()\.License(); pushServer == model\.MHPNS && (license == nil || !\*license\.Features\.MHPNS)/\/\/ Open source license always has MHPNS enabled/g' /Users/nghia/DEV/mattermost/server/channels/app/notification.go
sed -i '' 's/license == nil && maxUsersLimit > 0/false/g' /Users/nghia/DEV/mattermost/server/channels/app/limits.go
sed -i '' 's/license != nil && license\.IsSeatCountEnforced && license\.Features != nil && license\.Features\.Users != nil/true/g' /Users/nghia/DEV/mattermost/server/channels/app/limits.go
sed -i '' 's/license != nil && license\.Limits != nil && license\.Limits\.PostHistory > 0/true/g' /Users/nghia/DEV/mattermost/server/channels/app/limits.go
sed -i '' 's/license == nil || license\.Limits == nil || license\.Limits\.PostHistory == 0/false/g' /Users/nghia/DEV/mattermost/server/channels/app/limits.go
sed -i '' 's/license := a\.Srv()\.License(); license != nil && \*license\.Features\.LDAP/\/\/ Open source license always has LDAP enabled/g' /Users/nghia/DEV/mattermost/server/channels/app/ldap.go
sed -i '' 's/ldapI := a\.LdapDiagnostic(); ldapI != nil && license != nil && \*license\.Features\.LDAP && (\*a\.Config()\.LdapSettings\.Enable || \*a\.Config()\.LdapSettings\.EnableSync)/ldapI := a\.LdapDiagnostic(); ldapI != nil && (\*a\.Config()\.LdapSettings\.Enable || \*a\.Config()\.LdapSettings\.EnableSync)/g' /Users/nghia/DEV/mattermost/server/channels/app/ldap.go
sed -i '' 's/ldapI != nil && license != nil && model\.SafeDereference(license\.Features\.LDAP)/ldapI != nil/g' /Users/nghia/DEV/mattermost/server/channels/app/ldap.go
sed -i '' 's/ldapI != nil && license != nil && \*license\.Features\.LDAP/ldapI != nil/g' /Users/nghia/DEV/mattermost/server/channels/app/ldap.go
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/file.go
sed -i '' 's/license == nil || !license\.IsCloud()/false/g' /Users/nghia/DEV/mattermost/server/channels/app/file.go
sed -i '' 's/license := es\.license(); license != nil && \*license\.Features\.EmailNotificationContents/\/\/ Open source license always has EmailNotificationContents enabled/g' /Users/nghia/DEV/mattermost/server/channels/app/email/email_batching.go
sed -i '' 's/license != nil && \*license\.Features\.Compliance/true/g' /Users/nghia/DEV/mattermost/server/channels/app/email/email.go
sed -i '' 's/license := a\.Srv()\.License(); !\*a\.Config()\.ComplianceSettings\.Enable || license == nil || !\*license\.Features\.Compliance/\/\/ Open source license always has Compliance enabled\n\tif !\*a\.Config()\.ComplianceSettings\.Enable/g' /Users/nghia/DEV/mattermost/server/channels/app/compliance.go
sed -i '' 's/license := a\.Channels()\.License(); license == nil || !\*license\.Features\.MFA || !\*a\.Config()\.ServiceSettings\.EnableMultifactorAuthentication || !\*a\.Config()\.ServiceSettings\.EnforceMultifactorAuthentication/\/\/ Open source license always has MFA enabled\n\tif !\*a\.Config()\.ServiceSettings\.EnableMultifactorAuthentication || !\*a\.Config()\.ServiceSettings\.EnforceMultifactorAuthentication/g' /Users/nghia/DEV/mattermost/server/channels/app/authentication.go
sed -i '' 's/ldapAvailable := \*a\.Config()\.LdapSettings\.Enable && a\.Ldap() != nil && license != nil && \*license\.Features\.LDAP/ldapAvailable := \*a\.Config()\.LdapSettings\.Enable && a\.Ldap() != nil/g' /Users/nghia/DEV/mattermost/server/channels/app/authentication.go

# Remove license checks from API files
sed -i '' 's/license := c\.App\.Channels()\.License(); license == nil || !\*license\.Features\.CustomTermsOfService/\/\/ Open source license always has CustomTermsOfService enabled/g' /Users/nghia/DEV/mattermost/server/channels/api4/terms_of_service.go
sed -i '' 's/c\.App\.IPFiltering() == nil || !ipFilteringFeatureFlag || license == nil || !license\.IsCloud() || !model\.MinimumEnterpriseLicense(license)/c\.App\.IPFiltering() == nil || !ipFilteringFeatureFlag/g' /Users/nghia/DEV/mattermost/server/channels/api4/ip_filtering.go

# Remove license checks from web context
sed -i '' 's/license := c\.App\.Channels()\.License(); license == nil || !license\.IsCloud() || c\.AppContext\.Session()\.Props\[model\.SessionPropType\] != model\.SessionTypeCloudKey/\/\/ Open source license is not cloud\n\tif c\.AppContext\.Session()\.Props\[model\.SessionPropType\] != model\.SessionTypeCloudKey/g' /Users/nghia/DEV/mattermost/server/channels/web/context.go
sed -i '' 's/license := c\.App\.Channels()\.License(); license == nil || !license\.HasRemoteClusterService() || c\.AppContext\.Session()\.Props\[model\.SessionPropType\] != model\.SessionTypeRemoteclusterToken/\/\/ Open source license always has RemoteClusterService enabled\n\tif c\.AppContext\.Session()\.Props\[model\.SessionPropType\] != model\.SessionTypeRemoteclusterToken/g' /Users/nghia/DEV/mattermost/server/channels/web/context.go

echo "License checks removal completed!"
