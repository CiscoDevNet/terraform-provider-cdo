package license

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/sliceutil"
	"strconv"
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

var All = []Type{
	Base, Carrier, Threat, Malware, URLFilter,
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
		return fmt.Errorf("cannot unmarshal empty tring as a license type, it should be one of valid roles: %+v", nameToTypeMap)
	}
	unquoteType, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	deserialized, err := deserialize(unquoteType)
	if err != nil {
		return err
	}
	*t = deserialized
	return nil
}

func deserialize(name string) (Type, error) {
	l, ok := nameToTypeMap[name]
	if !ok {
		return "", fmt.Errorf("FTD License of name: \"%s\" not found, should be one of: %+v", name, nameToTypeMap)
	}
	return l, nil
}

func DeserializeAll(names string) ([]Type, error) {
	licenseStrs := strings.Split(names, ",")
	licenses := make([]Type, len(licenseStrs))
	for i, name := range licenseStrs {
		t, err := deserialize(name)
		if err != nil {
			return nil, err
		}
		licenses[i] = t
	}
	return licenses, nil
}

func SerializeAll(licenses []Type) string {
	return strings.Join(sliceutil.Map(licenses, func(l Type) string { return serialize(l) }), ",")
}

func serialize(license Type) string {
	return string(license)
}
