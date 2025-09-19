// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package app

import (
	"net/http"

	"github.com/mattermost/mattermost/server/public/model"
)

const (
	maxUsersLimit     = 200
	maxUsersHardLimit = 250
)

func (a *App) GetServerLimits() (*model.ServerLimits, *model.AppError) {
	limits := &model.ServerLimits{}
	license := a.License()

	// Determine user limits based on license presence and fields.
	if license == nil || license.Features == nil || license.Features.Users == nil {
		// Unlicensed or missing fields: apply hard-coded limits.
		limits.MaxUsersLimit = maxUsersLimit
		limits.MaxUsersHardLimit = maxUsersHardLimit
	} else {
		// Licensed: enforce license limits (no extra users unless configured elsewhere).
		licenseUserLimit := int64(*license.Features.Users)
		limits.MaxUsersLimit = licenseUserLimit
		// Extra users currently not configured here; default to 0.
		extraUsers := int64(0)
		limits.MaxUsersHardLimit = licenseUserLimit + extraUsers
	}

	// Apply post history limits only when license (and fields) are present.
	if license != nil && license.Limits != nil {
		limits.PostHistoryLimit = license.Limits.PostHistory
		// Get the calculated timestamp of the last accessible post
		lastAccessibleTime, appErr := a.GetLastAccessiblePostTime()
		if appErr != nil {
			return nil, appErr
		}
		limits.LastAccessiblePostTime = lastAccessibleTime
	}

	activeUserCount, appErr := a.Srv().Store().User().Count(model.UserCountOptions{})
	if appErr != nil {
		return nil, model.NewAppError("GetServerLimits", "app.limits.get_app_limits.user_count.store_error", nil, "", http.StatusInternalServerError).Wrap(appErr)
	}
	limits.ActiveUserCount = activeUserCount

	return limits, nil
}
func (a *App) GetPostHistoryLimit() int64 {
	license := a.License()
	if license == nil || license.Limits == nil {
		// No limits applicable
		return 0
	}

	return license.Limits.PostHistory
}

func (a *App) isAtUserLimit() (bool, *model.AppError) {
	userLimits, appErr := a.GetServerLimits()
	if appErr != nil {
		return false, appErr
	}

	if userLimits.MaxUsersHardLimit == 0 {
		return false, nil
	}

	return userLimits.ActiveUserCount >= userLimits.MaxUsersHardLimit, appErr
}
