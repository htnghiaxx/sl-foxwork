// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package einterfaces

import (
	"github.com/mattermost/mattermost/server/public/model"
)

// OpenSourceLicenseInterface defines the interface for open source license management
type OpenSourceLicenseInterface interface {
	GetLicense() *model.License
	IsLicensed() bool
	HasFeature(feature string) bool
	CanStartTrial() (bool, error)
	GetPrevTrial() (*model.License, error)
	NewMattermostEntryLicense(serverId string) *model.License
}

// OpenSourceLicenseManager implements the open source license interface
type OpenSourceLicenseManager struct {
	license *model.License
}

// NewOpenSourceLicenseManager creates a new open source license manager
func NewOpenSourceLicenseManager() *OpenSourceLicenseManager {
	return &OpenSourceLicenseManager{
		license: model.NewOpenSourceLicense(),
	}
}

// GetLicense returns the open source license
func (osm *OpenSourceLicenseManager) GetLicense() *model.License {
	return osm.license
}

// IsLicensed always returns true for open source
func (osm *OpenSourceLicenseManager) IsLicensed() bool {
	return true
}

// HasFeature always returns true for open source (all features enabled)
func (osm *OpenSourceLicenseManager) HasFeature(feature string) bool {
	return true
}

// CanStartTrial always returns false for open source (no trials needed)
func (osm *OpenSourceLicenseManager) CanStartTrial() (bool, error) {
	return false, nil
}

// GetPrevTrial returns nil for open source (no trials)
func (osm *OpenSourceLicenseManager) GetPrevTrial() (*model.License, error) {
	return nil, nil
}

// NewMattermostEntryLicense returns the open source license
func (osm *OpenSourceLicenseManager) NewMattermostEntryLicense(serverId string) *model.License {
	return osm.license
}
