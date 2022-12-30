// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import (
	"bytes"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateLicense(t *testing.T) {
	t.Run("should fail with junk data", func(t *testing.T) {
		b1 := []byte("junk")
		ok, _ := LicenseValidator.ValidateLicense(b1)
		require.False(t, ok, "should have failed - bad license")

		b2 := []byte("junkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunk")
		ok, _ = LicenseValidator.ValidateLicense(b2)
		require.False(t, ok, "should have failed - bad license")
	})

	t.Run("should not panic on shorted than expected input", func(t *testing.T) {
		var licenseData bytes.Buffer
		var inputData []byte

		for i := 0; i < 255; i++ {
			inputData = append(inputData, 'A')
		}
		inputData = append(inputData, 0x00)

		encoder := base64.NewEncoder(base64.StdEncoding, &licenseData)
		_, err := encoder.Write(inputData)
		require.NoError(t, err)
		err = encoder.Close()
		require.NoError(t, err)

		ok, str := LicenseValidator.ValidateLicense(licenseData.Bytes())
		require.False(t, ok)
		require.Empty(t, str)
	})

	t.Run("should not panic with input filled of null terminators", func(t *testing.T) {
		var licenseData bytes.Buffer
		var inputData []byte

		for i := 0; i < 256; i++ {
			inputData = append(inputData, 0x00)
		}

		encoder := base64.NewEncoder(base64.StdEncoding, &licenseData)
		_, err := encoder.Write(inputData)
		require.NoError(t, err)
		err = encoder.Close()
		require.NoError(t, err)

		ok, str := LicenseValidator.ValidateLicense(licenseData.Bytes())
		require.False(t, ok)
		require.Empty(t, str)
	})
}

func TestGetLicenseFileLocation(t *testing.T) {
	fileName := GetLicenseFileLocation("")
	require.NotEmpty(t, fileName, "invalid default file name")

	fileName = GetLicenseFileLocation("mattermost.mattermost-license")
	require.Equal(t, fileName, "mattermost.mattermost-license", "invalid file name")
}

func TestGetLicenseFileFromDisk(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		fileBytes := GetLicenseFileFromDisk("thisfileshouldnotexist.mattermost-license")
		assert.Empty(t, fileBytes, "invalid bytes")
	})

	t.Run("not a license file", func(t *testing.T) {
		f, err := os.CreateTemp("", "TestGetLicenseFileFromDisk")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		os.WriteFile(f.Name(), []byte("not a license"), 0777)

		fileBytes := GetLicenseFileFromDisk(f.Name())
		require.NotEmpty(t, fileBytes, "should have read the file")

		success, _ := LicenseValidator.ValidateLicense(fileBytes)
		assert.False(t, success, "should have been an invalid file")
	})
}

func TestGetNextTrueUpReviewDueDate(t *testing.T) {
	t.Run("Due date always falls on the 15th", func(t *testing.T) {
		// Before the 15th
		now := time.Date(2022, 12, 14, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, due.Day(), trueUpReviewDueDay)

		// On the 15th
		now = time.Date(2022, 12, 15, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, due.Day(), trueUpReviewDueDay)

		// After the 15th
		now = time.Date(2022, 12, 16, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, due.Day(), trueUpReviewDueDay)
	})

	t.Run("Due date will always be in next quarter if the current date is past the 15th", func(t *testing.T) {
		now := time.Date(2022, time.March, 16, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.June, due.Month())

		now = time.Date(2022, time.June, 16, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.September, due.Month())

		now = time.Date(2022, time.September, 16, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.December, due.Month())

		now = time.Date(2022, time.December, 16, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.March, due.Month())
	})

	t.Run("Due date will always be in the current quarter if the current date is before or on the 15th", func(t *testing.T) {
		now := time.Date(2022, time.March, 15, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.March, due.Month())

		now = time.Date(2022, time.June, 15, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.June, due.Month())

		now = time.Date(2022, time.September, 14, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.September, due.Month())

		now = time.Date(2022, time.December, 14, 0, 0, 0, 0, time.Local)
		due = GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.December, due.Month())
	})

	t.Run("Due date will be in the next year if the next quarter is not within the current year", func(t *testing.T) {
		now := time.Date(2022, time.December, 18, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(now)
		assert.Equal(t, time.March, due.Month())
		assert.Equal(t, 2023, due.Year())
	})
}

func TestIsTrueUpReviewDueDateWithinTheNextTwoWeeks(t *testing.T) {
	t.Run("Ensure a date within two weeks before the due date returns true", func(t *testing.T) {
		// 1 Day before the due date
		now := time.Date(2022, time.December, 14, 0, 0, 0, 0, time.Local)
		// Due date is December 15th, 2022
		due := GetNextTrueUpReviewDueDate(now)

		res := IsTrueUpReviewDueDateWithinTheNextTwoWeeks(now, due)
		assert.True(t, res)
	})

	t.Run("Ensure a date that is more than two weeks before the due date returns false", func(t *testing.T) {
		// 15 Days before the due date
		now := time.Date(2022, time.November, 30, 0, 0, 0, 0, time.Local)
		// Due date is December 15th, 2022
		due := GetNextTrueUpReviewDueDate(now)

		res := IsTrueUpReviewDueDateWithinTheNextTwoWeeks(now, due)
		assert.False(t, res)
	})

	t.Run("Ensure a date that past the due date returns false", func(t *testing.T) {
		now := time.Date(2022, time.December, 16, 0, 0, 0, 0, time.Local)

		// Due date is December 15th, 2022
		dueNow := time.Date(2022, time.December, 15, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(dueNow)

		res := IsTrueUpReviewDueDateWithinTheNextTwoWeeks(now, due)
		assert.False(t, res)
	})

	t.Run("Ensure a date that is on the due date returns true", func(t *testing.T) {
		now := time.Date(2022, time.December, 15, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(now)

		res := IsTrueUpReviewDueDateWithinTheNextTwoWeeks(now, due)
		assert.True(t, res)
	})

	t.Run("Ensure a date that is on the first day of the due date window returns true", func(t *testing.T) {
		now := time.Date(2022, time.December, 1, 0, 0, 0, 0, time.Local)
		due := GetNextTrueUpReviewDueDate(now)

		res := IsTrueUpReviewDueDateWithinTheNextTwoWeeks(now, due)
		assert.True(t, res)
	})
}
