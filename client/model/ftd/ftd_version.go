package ftd

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// The compiled version of the regex created at init() is cached here so it
// only needs to be created once.
var versionRegex *regexp.Regexp

// semVerRegex is the regular expression used to parse a semantic version.
// This is not the official regex from the semver spec. It has been modified to allow for loose handling
// where versions like 2.1 are detected.
const semVerRegex = `^(\d+)\.(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:-(\d+))?$`

// Version represents a single semantic version.
type Version struct {
	major, minor, patch, hotfix uint64
	buildNum                    uint64
	original                    string
}

func init() {
	versionRegex = regexp.MustCompile("^" + semVerRegex + "$")
}

const (
	num     string = "0123456789"
	allowed string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-" + num
)

// NewVersion parses a given version and returns an instance of Version or
// an error if unable to parse the version. If the version is SemVer-ish it
// attempts to convert it to SemVer. If you want  to validate it was a strict
// semantic version at parse time see StrictNewVersion().
func NewVersion(v string) (*Version, error) {
	m := versionRegex.FindStringSubmatch(v)
	if m == nil {
		return nil, errors.New(fmt.Sprintf("invalid FTD Version %s", v))
	}

	sv := &Version{
		original: v,
	}

	var err error
	sv.major, err = strconv.ParseUint(m[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("Error parsing version segment: %s", err)
	}

	if m[2] != "" {
		sv.minor, err = strconv.ParseUint(m[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
	} else {
		sv.minor = 0
	}

	if m[3] != "" {
		sv.patch, err = strconv.ParseUint(m[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
	} else {
		sv.patch = 0
	}

	if m[4] != "" {
		sv.hotfix, err = strconv.ParseUint(m[4], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
	} else {
		sv.hotfix = 0
	}

	if m[5] != "" {
		sv.buildNum, err = strconv.ParseUint(m[5], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing version segment: %s", err)
		}
	} else {
		sv.buildNum = 0
	}

	return sv, nil
}

// String converts a Version object to a string.
// Note, if the original version contained a leading v this version will not.
// See the Original() method to retrieve the original value. Semantic Versions
// don't contain a leading v per the spec. Instead it's optional on
// implementation.
func (v Version) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%d.%d.%d", v.major, v.minor, v.patch)
	if v.hotfix != 0 {
		fmt.Fprintf(&buf, ".%d", v.hotfix)
	}
	if v.buildNum != 0 {
		fmt.Fprintf(&buf, "-%d", v.buildNum)
	}

	return buf.String()
}

// Original returns the original value passed in to be parsed.
func (v *Version) Original() string {
	return v.original
}

// Major returns the major version.
func (v Version) Major() uint64 {
	return v.major
}

// Minor returns the minor version.
func (v Version) Minor() uint64 {
	return v.minor
}

// Patch returns the patch version.
func (v Version) Patch() uint64 {
	return v.patch
}

func (v Version) Hotfix() uint64 {
	return v.hotfix
}

// Prerelease returns the pre-release version.
func (v Version) Buildnum() uint64 {
	return v.buildNum
}

// LessThan tests if one version is less than another one.
func (v *Version) LessThan(o *Version) bool {
	return v.Compare(o) < 0
}

// LessThanEqual tests if one version is less or equal than another one.
func (v *Version) LessThanEqual(o *Version) bool {
	return v.Compare(o) <= 0
}

// GreaterThan tests if one version is greater than another one.
func (v *Version) GreaterThan(o *Version) bool {
	return v.Compare(o) > 0
}

// GreaterThanEqual tests if one version is greater or equal than another one.
func (v *Version) GreaterThanEqual(o *Version) bool {
	return v.Compare(o) >= 0
}

// Equal tests if two versions are equal to each other.
// Note, versions can be equal with different metadata since metadata
// is not considered part of the comparable version.
func (v *Version) Equal(o *Version) bool {
	if v == o {
		return true
	}
	if v == nil || o == nil {
		return false
	}
	return v.Compare(o) == 0
}

// Compare compares this version to another one. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version.
//
// Versions are compared by X.Y.Z. Build metadata is ignored. Prerelease is
// lower than the version without a prerelease. Compare always takes into account
// prereleases. If you want to work with ranges using typical range syntaxes that
// skip prereleases if the range is not looking for them use constraints.
func (v *Version) Compare(o *Version) int {
	if d := compareSegment(v.Major(), o.Major()); d != 0 {
		return d
	}
	if d := compareSegment(v.Minor(), o.Minor()); d != 0 {
		return d
	}
	if d := compareSegment(v.Patch(), o.Patch()); d != 0 {
		return d
	}
	if d := compareSegment(v.Hotfix(), o.Hotfix()); d != 0 {
		return d
	}
	if d := compareSegment(v.Buildnum(), o.Buildnum()); d != 0 {
		return d
	}

	return 0
}

func compareSegment(v, o uint64) int {
	if v < o {
		return -1
	}
	if v > o {
		return 1
	}

	return 0
}
