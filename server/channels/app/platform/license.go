// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

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
	// Open source license is always available, no need to load anything
	ps.logger.Info("Open source license loaded, all features enabled.")
}

func (ps *PlatformService) SetLicense(license *model.License) bool {
	// Always return true for open source, license is always available
	return true
}

func (ps *PlatformService) ClientLicense() map[string]string {
	return map[string]string{
		"IsLicensed":                "true",
		"Users":                     "999999999",
		"LDAP":                      "true",
		"LDAPGroups":                "true",
		"MFA":                       "true",
		"GoogleOAuth":               "true",
		"Office365OAuth":            "true",
		"OpenId":                    "true",
		"Compliance":                "true",
		"Cluster":                   "true",
		"Metrics":                   "true",
		"MHPNS":                     "true",
		"SAML":                      "true",
		"Elasticsearch":             "true",
		"Announcement":              "true",
		"ThemeManagement":           "true",
		"EmailNotificationContents": "true",
		"DataRetention":             "true",
		"MessageExport":             "true",
		"CustomPermissionsSchemes":  "true",
		"CustomTermsOfService":      "true",
		"GuestAccounts":             "true",
		"GuestAccountsPermissions":  "true",
		"IDLoadedPushNotifications": "true",
		"LockTeammateNameDisplay":   "true",
		"EnterprisePlugins":         "true",
		"AdvancedLogging":           "true",
		"Cloud":                     "false",
		"SharedChannels":            "true",
		"RemoteClusterService":      "true",
		"OutgoingOAuthConnections":  "true",
		"FutureFeatures":            "true",
	}
}

func (ps *PlatformService) GetSanitizedClientLicense() map[string]string {
	return ps.ClientLicense()
}

func (ps *PlatformService) AddLicenseListener(listener func(oldLicense, newLicense *model.License)) string {
	// For open source, we don't need license listeners since license never changes
	return "opensource-listener"
}

func (ps *PlatformService) RemoveLicenseListener(id string) {
	// No-op for open source
}

func (ps *PlatformService) RemoveLicense() *model.AppError {
	// No-op for open source, license is always available
	return nil
}

func (ps *PlatformService) ValidateAndSetLicenseBytes(b []byte) error {
	// No-op for open source, license is always valid
	return nil
}

func (ps *PlatformService) SetClientLicense(m map[string]string) {
	// No-op for open source, client license is always the same
}

func (ps *PlatformService) SaveLicense(licenseBytes []byte) (*model.License, *model.AppError) {
	// Return the open source license
	return ps.licenseManager.GetLicense(), nil
}

func (ps *PlatformService) RequestTrialLicense(trialRequest *model.TrialLicenseRequest) *model.AppError {
	// No-op for open source, no trials needed
	return nil
}
