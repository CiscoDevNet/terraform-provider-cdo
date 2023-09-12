package role

import (
	"fmt"
	"strconv"
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

var All = []Type{
	User,
	Admin,
	ReadOnly,
	SuperAdmin,
	DeployOnly,
	EditOnly,
	VpnSessionsManager,
	ChoiceAdmin,
}

var nameToType = make(map[string]Type, len(All))

func init() {
	for _, t := range All {
		nameToType[string(t)] = t
	}
}

func (t *Type) UnmarshalJSON(b []byte) error {
	if len(b) <= 2 || b == nil {
		return fmt.Errorf("cannot unmarshal empty tring as a role type, it should be one of valid roles: %+v", nameToType)
	}
	unquotedType, err := strconv.Unquote(string(b))
	if err != nil {
		return err
	}
	role, ok := nameToType[unquotedType]
	if !ok {
		return fmt.Errorf("cannot unmarshal %s as a role type, it should be one of valid roles: %+v", string(b), nameToType)
	}
	*t = role
	return nil
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(string(*t))), nil
}
