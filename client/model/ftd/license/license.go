package license

import (
	"fmt"
	"strings"
)

type Type string

//type TypeSlice []Type

//func (t *Type) UnmarshalJSON(data []byte) error {
//	if len(data) == 0 || string(data) == "null" {
//		return nil
//	}
//	dataStr := string(data)
//	for _, l := range licenseMap {
//		if string(l) == dataStr {
//			*t = l
//		}
//	}
//	return fmt.Errorf("cannot unmarshal json: \"%s\" to type license.Type, it should be one of %s", string(data), licenseMap)
//}
//
//func (ts *TypeSlice) UnmarshalJSON(data []byte) error {
//	if len(data) == 0 || string(data) == "null" {
//		return nil
//	}
//	dataStr := string(data)
//	licenseStrs := strings.Split(dataStr, ",")
//	licenseSlice := make([]Type, len(licenseStrs))
//	for i, licenseStr := range licenseStrs {
//		license, ok := licenseMap[licenseStr]
//		if !ok {
//			return fmt.Errorf("cannot unmarshal json, license: \"%s\" is not a valid license, valid licenses: %+v", license, licenseMap)
//		}
//		licenseSlice[i] = license
//	}
//
//	*ts = licenseSlice
//	return nil
//}

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

func Parse(name string) (Type, error) {
	l, ok := licenseMap[name]
	if !ok {
		return "", fmt.Errorf("FTD License of name: \"%s\" not found, should be one of: %+v", name, licenseMap)
	}
	return l, nil
}

func ParseAll(names string) ([]Type, error) {
	licenseStrs := strings.Split(names, ",")
	licenses := make([]Type, len(licenseStrs))
	for i, name := range licenseStrs {
		t, err := Parse(name)
		if err != nil {
			return nil, err
		}
		licenses[i] = t
	}
	return licenses, nil
}
