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
	Base,
	Carrier,
	Threat,
	Malware,
	URLFilter,
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
	deserialized, err := Deserialize(string(b[1 : len(b)-1])) // strip off quote
	if err != nil {
		return err
	}
	*t = deserialized
	return nil
}

func MustParse(name string) Type {
	l, ok := nameToTypeMap[name]
	if !ok {
		panic(fmt.Errorf("FTD License of name: \"%s\" not found, should be one of %+v", name, AllAsString))
	}
	return l
}

func Deserialize(name string) (Type, error) {
	l, ok := nameToTypeMap[name]
	if !ok {
		return "", fmt.Errorf("FTD License of name: \"%s\" not found, should be one of: %+v", name, AllAsString)
	}
	return l, nil
}

func SerializeAllAsCdo(licenses []Type) string {
	return strings.Join(sliceutil.Map(licenses, func(l Type) string { return string(l) }), ",")
}

// DeserializeAllFromCdo exists because CDO store license caps as one comma-sep string
// but fmc store it as list of string, use this method to handle CDO's special case
func DeserializeAllFromCdo(licenses string) ([]Type, error) {
	return sliceutil.MapWithError(strings.Split(licenses, ","), func(l string) (Type, error) {
		t, ok := nameToTypeMap[l]
		if !ok {
			return "", fmt.Errorf("cannot deserialize %s as license, should be one of %+v", l, All)
		}
		return t, nil
	})
}
