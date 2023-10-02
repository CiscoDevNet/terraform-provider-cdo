package status

type Type string

// see also: https://github.com/cisco-lockhart/lh-plugin-core/blob/master/lh-target/src/main/java/com/cisco/lockhart/target/data/ServiceActivationState.java
const (
	New         Type = "NEW"
	ReOnboard   Type = "REONBOARD"
	YoReOnboard Type = "YO_REONBOARD"
	Onboarding  Type = "ONBOARDING"
	Active      Type = "ACTIVE"
	Inactive    Type = "INACTIVE"
	Disabled    Type = "DISABLE"
)
