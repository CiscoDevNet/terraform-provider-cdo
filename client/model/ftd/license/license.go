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

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte("\"" + Serialize(*t) + "\""), nil
}

func (t *Type) UnmarshalJSON(b []byte) error {
	if len(b) <= 2 || b == nil {
		return fmt.Errorf("cannot unmarshal empty tring as a license type, it should be one of valid roles: %+v", licenseMap)
	}
	deserialized, err := Deserialize(string(b[1 : len(b)-1])) // strip off quote
	if err != nil {
		return err
	}
	*t = deserialized
	return nil
}

//func (t *Types) MarshalJSON() ([]byte, error) {
//	return []byte("\"" + SerializeAll(*t) + "\""), nil
//}
//
//func (t *Types) UnmarshalJSON(b []byte) error {
//	if len(b) <= 2 || b == nil {
//		return fmt.Errorf("cannot unmarshal empty tring as a license type, it should be one of valid roles: %+v", licenseMap)
//	}
//	deserialized, err := DeserializeAll(string(b[1 : len(b)-1])) // strip off bracket
//	if err != nil {
//		return err
//	}
//	*t = deserialized
//	return nil
//}

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
		panic(fmt.Errorf("FTD License of name: \"%s\" not found, should be one of %+v", name, licenseMap))
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
