package role

import "fmt"

type Type string

const (
	USER               Type = "ROLE_USER"
	ADMIN              Type = "ROLE_ADMIN"
	ReadOnly           Type = "ROLE_READ_ONLY"
	SuperAdmin         Type = "ROLE_SUPER_ADMIN"
	DeployOnly         Type = "ROLE_DEPLOY_ONLY"
	EditOnly           Type = "ROLE_EDIT_ONLY"
	VpnSessionsManager Type = "ROLE_VPN_SESSIONS_MANAGER"
	ChoiceAdmin        Type = "ROLE_CHOICE_ADMIN"
)

var NameToType = map[string]Type{
	"ROLE_USER":                 USER,
	"ROLE_ADMIN":                ADMIN,
	"ROLE_READ_ONLY":            ReadOnly,
	"ROLE_SUPER_ADMIN":          SuperAdmin,
	"ROLE_DEPLOY_ONLY":          DeployOnly,
	"ROLE_EDIT_ONLY":            EditOnly,
	"ROLE_VPN_SESSIONS_MANAGER": VpnSessionsManager,
	"ROLE_CHOICE_ADMIN":         ChoiceAdmin,
}

func (t *Type) UnmarshalJSON(b []byte) error {
	role, ok := NameToType[string(b)]
	if !ok {
		return fmt.Errorf("cannot unmarshal %s as a role type, it should be one of valid roles: %+v", string(b), NameToType)
	}
	*t = role
	return nil
}
