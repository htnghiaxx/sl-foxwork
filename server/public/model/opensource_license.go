// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package model

// OpenSourceLicense represents a permanent open source license with all features enabled
type OpenSourceLicense struct {
	Id        string    `json:"id"`
	IssuedAt  int64     `json:"issued_at"`
	StartsAt  int64     `json:"starts_at"`
	ExpiresAt int64     `json:"expires_at"`
	Features  *Features `json:"features"`
}

// NewOpenSourceLicense creates a new open source license with all features enabled
func NewOpenSourceLicense() *License {
	now := GetMillis()
	// Set expiration to 100 years from now
	expiresAt := now + (100 * 365 * 24 * 60 * 60 * 1000)

	return &License{
		Id:           "opensource-permanent",
		IssuedAt:     now,
		StartsAt:     now,
		ExpiresAt:    expiresAt,
		SkuName:      "Open Source",
		SkuShortName: "opensource",
		IsTrial:      false,
		IsGovSku:     false,
		Features: &Features{
			Users:                     NewPointer(999999999), // Unlimited users
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

// IsOpenSourceLicense checks if a license is an open source license
func (l *License) IsOpenSourceLicense() bool {
	return l != nil && l.Id == "opensource-permanent"
}
