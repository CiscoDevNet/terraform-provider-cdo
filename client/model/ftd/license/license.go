package license

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/sliceutil"
	"strings"
)

type Type string

// https://www.cisco.com/c/en/us/td/docs/security/firepower/70/fdm/fptd-fdm-config-guide-700/fptd-fdm-license.html
const (
	Base      Type = "BASE"
	Carrier   Type = "CARRIER"
	Threat    Type = "THREAT"
	Malware   Type = "MALWARE"
	URLFilter Type = "URLFilter"
)

var licenseMap = map[string]Type{
	"BASE":      Base,
	"CARRIER":   Carrier,
	"THREAT":    Threat,
	"MALWARE":   Malware,
	"URLFilter": URLFilter,
}

func MustParse(name string) Type {
	l, ok := licenseMap[name]
	if !ok {
		panic(fmt.Errorf("FTD License of name: \"%s\" not found", name))
	}
	return l
}

func Deserialize(name string) (Type, error) {
	l, ok := licenseMap[name]
	if !ok {
		return "", fmt.Errorf("FTD License of name: \"%s\" not found, should be one of: %+v", name, licenseMap)
	}
	return l, nil
}

func DeserializeAll(names string) ([]Type, error) {
	licenseStrs := strings.Split(names, ",")
	licenses := make([]Type, len(licenseStrs))
	for i, name := range licenseStrs {
		t, err := Deserialize(name)
		if err != nil {
			return nil, err
		}
		licenses[i] = t
	}
	return licenses, nil
}

func SerializeAll(licenses []Type) string {
	return strings.Join(sliceutil.Map(licenses, func(l Type) string { return Serialize(l) }), ",")
}

func Serialize(license Type) string {
	return string(license)
}
