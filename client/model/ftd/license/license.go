package license

import "fmt"

type Type string

// https://www.cisco.com/c/en/us/td/docs/security/firepower/70/fdm/fptd-fdm-config-guide-700/fptd-fdm-license.html
const (
	Base      Type = "BASE"
	Carrier   Type = "CARRIER"
	Threat    Type = "THREAT"
	Malware   Type = "MALWARE"
	URLFilter Type = "URLFilter"
)

var AllLicenses = []Type{
	Base,
	Carrier,
	Threat,
	Malware,
	URLFilter,
}

func MustParse(name string) Type {
	for _, l := range AllLicenses {
		if string(l) == name {
			return l
		}
	}
	panic(fmt.Errorf("FTD License of name: \"%s\" not found", name))
}
