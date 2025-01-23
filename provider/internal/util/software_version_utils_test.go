package util_test

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"testing"
)

func TestNormalizeAsaVersion(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        string
		expResult string
	}

	testCases := map[string]testCase{
		"normal": {
			in:        "9.8(1)",
			expResult: "9.8.1",
		},
		"no-parentheses": {
			in:        "9.8.1",
			expResult: "9.8.1",
		},
		"multiple-parentheses": {
			in:        "9.8(1)(201)",
			expResult: "9.8.1.201",
		},
		"parentheses": {
			in:        "9.8(1.201)",
			expResult: "9.8.1.201",
		},
		"single-parentheses": {
			in:        "9.8(1)21",
			expResult: "9.8.1.21",
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := util.NormalizeAsaVersion(test.in)
			if result != test.expResult {
				t.Errorf("expected %s, got %s", test.expResult, result)
			}
		})
	}
}

func TestDoNormalisedVersionsMatch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		versionOne string
		versionTwo string
		expResult  bool
	}

	testCases := map[string]testCase{
		"matches-paren": {
			versionOne: "9.8(1)",
			versionTwo: "9.8.1",
			expResult:  true,
		},
		"matches-paren-with-dots": {
			versionOne: "9.8(1.100)",
			versionTwo: "9.8(1)100",
			expResult:  true,
		},
		"no-parentheses-with-dots": {
			versionOne: "9.8(1.100)",
			versionTwo: "9.8(1.101)",
			expResult:  true,
		},
	}

	for name, test := range testCases {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			result := util.DoNormalisedVersionsMatch(test.versionOne, test.versionTwo)
			if result != test.expResult {
				t.Errorf("expected %t, got %t", test.expResult, result)
			}
		})
	}
}
