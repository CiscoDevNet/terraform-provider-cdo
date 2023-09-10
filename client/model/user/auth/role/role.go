package role

import (
	"fmt"
)

type Type string

const (
	User               Type = "ROLE_USER"
	Admin              Type = "ROLE_ADMIN"
	ReadOnly           Type = "ROLE_READ_ONLY"
	SuperAdmin         Type = "ROLE_SUPER_ADMIN"
	DeployOnly         Type = "ROLE_DEPLOY_ONLY"
	EditOnly           Type = "ROLE_EDIT_ONLY"
	VpnSessionsManager Type = "ROLE_VPN_SESSIONS_MANAGER"
	ChoiceAdmin        Type = "ROLE_CHOICE_ADMIN"
)

var NameToType = map[string]Type{
	string(User):               User,
	string(Admin):              Admin,
	string(ReadOnly):           ReadOnly,
	string(SuperAdmin):         SuperAdmin,
	string(DeployOnly):         DeployOnly,
	string(EditOnly):           EditOnly,
	string(VpnSessionsManager): VpnSessionsManager,
	string(ChoiceAdmin):        ChoiceAdmin,
}

func (t *Type) UnmarshalJSON(b []byte) error {
	if len(b) <= 2 || b == nil {
		return fmt.Errorf("cannot unmarshal empty tring as a role type, it should be one of valid roles: %+v", NameToType)
	}
	role, ok := NameToType[string(b[1:len(b)-1])]
	if !ok {
		return fmt.Errorf("cannot unmarshal %s as a role type, it should be one of valid roles: %+v", string(b), NameToType)
	}
	*t = role
	return nil
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(*t), nil
}
