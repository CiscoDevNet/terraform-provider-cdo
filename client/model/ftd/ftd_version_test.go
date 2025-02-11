package ftd_test

import (
	"errors"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewFtdVersion(t *testing.T) {
	testCases := []struct {
		testName    string
		versionStr  string
		expectedErr error
	}{
		{
			testName:    "Successfully parse major.minor.patch",
			versionStr:  "7.2.0",
			expectedErr: nil,
		},
		{
			testName:    "Successfully parse major.minor.patch-buildnum",
			versionStr:  "7.2.0-69",
			expectedErr: nil,
		},
		{
			testName:    "Successfully parse major.minor.patch.hotfix",
			versionStr:  "7.2.1.45",
			expectedErr: nil,
		},
		{
			testName:    "Successfully parse major.minor.patch.hotfix-buildnumber",
			versionStr:  "7.2.1.45-59",
			expectedErr: nil,
		},
		{
			testName:    "Fail to parse invalid FTD version",
			versionStr:  "9.8.4(100)1",
			expectedErr: errors.New("invalid FTD Version 9.8.4(100)1"),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			version, err := ftd.NewVersion(testCase.versionStr)
			if err == nil {
				assert.Equal(t, testCase.versionStr, version.String())
			}
			assert.Equal(t, testCase.expectedErr, err)
		})
	}
}

func TestFtdVersionComparison(t *testing.T) {
	testCases := []struct {
		testName      string
		versionOneStr string
		versionTwoStr string
		assertFunc    func(t *testing.T, versionOne, versionTwo *ftd.Version)
	}{
		{
			testName:      "Compare major.minor.patch versions",
			versionOneStr: "7.2.0",
			versionTwoStr: "7.3.0",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.True(t, versionOne.LessThan(versionTwo))
				assert.True(t, versionTwo.GreaterThan(versionOne))
				assert.False(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare major.minor.patch versions semantically",
			versionOneStr: "7.1.2",
			versionTwoStr: "7.3.0",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.True(t, versionOne.LessThan(versionTwo))
				assert.True(t, versionTwo.GreaterThan(versionOne))
				assert.False(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare more major.minor.patch versions semantically",
			versionOneStr: "7.8.0",
			versionTwoStr: "7.12.0",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.True(t, versionOne.LessThan(versionTwo))
				assert.True(t, versionTwo.GreaterThan(versionOne))
				assert.False(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare equal major.minor.patch versions",
			versionOneStr: "7.2.0",
			versionTwoStr: "7.2.0",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.False(t, versionOne.LessThan(versionTwo))
				assert.False(t, versionTwo.LessThan(versionOne))
				assert.True(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare major.minor.patch-buildnum versions",
			versionOneStr: "7.2.0-68",
			versionTwoStr: "7.2.0-69",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.False(t, versionOne.LessThan(versionTwo))
				assert.False(t, versionTwo.GreaterThan(versionOne))
				assert.True(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare equal major.minor.patch-buildnum versions",
			versionOneStr: "7.2.0-65",
			versionTwoStr: "7.2.0-65",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.False(t, versionOne.LessThan(versionTwo))
				assert.False(t, versionTwo.LessThan(versionOne))
				assert.True(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare major.minor.patch-buildnum versions",
			versionOneStr: "7.2.0.2-68",
			versionTwoStr: "7.2.0.3-69",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.True(t, versionOne.LessThan(versionTwo))
				assert.True(t, versionTwo.GreaterThan(versionOne))
				assert.False(t, versionTwo.Equal(versionOne))
			},
		},
		{
			testName:      "Compare equal major.minor.patch.hotfix-buildnum versions",
			versionOneStr: "7.2.1.24-65",
			versionTwoStr: "7.2.1.24-65",
			assertFunc: func(t *testing.T, versionOne, versionTwo *ftd.Version) {
				assert.False(t, versionOne.LessThan(versionTwo))
				assert.False(t, versionTwo.LessThan(versionOne))
				assert.True(t, versionTwo.Equal(versionOne))
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			parsedVersionOne, err := ftd.NewVersion(testCase.versionOneStr)
			assert.Nil(t, err)
			parsedVersionTwo, err := ftd.NewVersion(testCase.versionTwoStr)
			assert.Nil(t, err)
			testCase.assertFunc(t, parsedVersionOne, parsedVersionTwo)
		})
	}
}
