package tier

import "fmt"

type Type string

// https://www.cisco.com/c/en/us/td/docs/security/firepower/70/fdm/fptd-fdm-config-guide-700/fptd-fdm-license.html
const (
	FTDv5   Type = "FTDv5"
	FTDv10  Type = "FTDv10"
	FTDv20  Type = "FTDv20"
	FTDv30  Type = "FTDv30"
	FTDv50  Type = "FTDv50"
	FTDv100 Type = "FTDv100"
	FTDv    Type = "FTDv"
)

var AllTiers = []Type{
	FTDv5,
	FTDv10,
	FTDv20,
	FTDv30,
	FTDv50,
	FTDv100,
	FTDv,
}

func Parse(name string) (Type, error) {
	for _, tier := range AllTiers {
		if string(tier) == name {
			return tier, nil
		}
	}
	return "", fmt.Errorf("FTD Performance Tier of name: \"%s\" not found", name)
}
