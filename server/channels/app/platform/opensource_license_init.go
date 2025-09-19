// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package platform

import (
	"github.com/mattermost/mattermost/server/v8/einterfaces"
)

func init() {
	// Register the open source license manager
	RegisterLicenseInterface(func(ps *PlatformService) einterfaces.LicenseInterface {
		return einterfaces.NewOpenSourceLicenseManager()
	})
}
