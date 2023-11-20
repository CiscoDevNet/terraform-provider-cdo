package license

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/sliceutil"
)

type Type string

// https://www.cisco.com/c/en/us/td/docs/security/firepower/70/fdm/fptd-fdm-config-guide-700/fptd-fdm-license.html
const (
	// CDO terms
	Base      Type = "BASE"
	Carrier   Type = "CARRIER"
	Threat    Type = "THREAT"
	Malware   Type = "MALWARE"
	URLFilter Type = "URLFilter"

	// FMC terms
	Essentials     Type = "ESSENTIALS"
	URL            Type = "URL"
	IPS            Type = "IPS"
	MalwareDefense Type = "MALWARE_DEFENSE"
)

var All = []Type{
	Base,
	Essentials,
	Carrier,
	Threat,
	IPS,
	Malware,
	MalwareDefense,
	URLFilter,
	URL,
}

var fmcToCdoLicenseNameMap = map[Type]Type{
	Essentials:     Base,
	IPS:            Threat,
	URL:            URLFilter,
	MalwareDefense: Malware,
}

var cdoToFmcLicenseNameMap = map[Type]Type{
	Base:      Essentials,
	Threat:    IPS,
	URLFilter: URL,
	Malware:   MalwareDefense,
}

var AllAsString = make([]string, len(All))

var nameToTypeMap = make(map[string]Type, len(All))

func init() {
	for i, l := range All {
		nameToTypeMap[string(l)] = l
		AllAsString[i] = string(l)
	}
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(*t))), nil
}

func (t *Type) UnmarshalJSON(b []byte) error {
	if len(b) <= 2 || b == nil {
		return fmt.Errorf("cannot unmarshal empty tring as a license type, it should be one of valid roles: %+v", AllAsString)
	}
	deserialized, err := stringToLicense(string(b[1 : len(b)-1])) // strip off quote
	if err != nil {
		return err
	}
	*t = deserialized
	return nil
}

// replaceFmcLicenseTermsWithCdoTerms is used to tell terraform how the licenses returned by FMC map to licenses expected by CDO
// We need to tell Terraform during read that they are the same thing, so when reading it back, we need the conversion
func replaceFmcTermWithCdoTerm(fmcLicense Type) Type {
	cdoLicense, ok := fmcToCdoLicenseNameMap[fmcLicense]
	if ok {
		return cdoLicense
	} else {
		return fmcLicense
	}
}

func replaceCdoTermWithFmcTerm(cdoLicense Type) Type {
	fmcLicense, ok := cdoToFmcLicenseNameMap[cdoLicense]
	if ok {
		return fmcLicense
	} else {
		return cdoLicense
	}
}

func stringToLicense(name string) (Type, error) {
	l, ok := nameToTypeMap[name]
	if !ok {
		return "", fmt.Errorf("FTD License of name: \"%s\" not found, should be one of: %+v", name, AllAsString)
	}
	return l, nil
}

func LicensesToString(licenses []Type) string {
	return strings.Join(LicensesToStrings(licenses), ",")
}

func LicensesToStrings(licenses []Type) []string {
	return sliceutil.Map(licenses, func(l Type) string { return string(l) })
}

// StringToCdoLicenses exists because CDO store license caps as one comma-sep string
// but fmc store it as list of string, and they have different name for some licenses,
// use this method to handle this special case by converting FMC representation to
// CDO representation
func StringToCdoLicenses(licenses string) ([]Type, error) {
	return sliceutil.MapWithError(strings.Split(licenses, ","), func(l string) (Type, error) {
		t, ok := nameToTypeMap[l]
		if !ok {
			return "", fmt.Errorf("cannot deserialize %s as license, should be one of %+v", l, All)
		}
		t = replaceFmcTermWithCdoTerm(t)
		return t, nil
	})
}

func LicensesToFmcLicenses(licenses []Type) []Type {
	return sliceutil.Map(licenses, func(l Type) Type {
		return replaceCdoTermWithFmcTerm(l)
	})
}

func StringToCdoStrings(licenses string) ([]string, error) {
	licenseTypes, err := StringToCdoLicenses(licenses)
	if err != nil {
		return nil, err
	}
	return LicensesToStrings(licenseTypes), nil
}
