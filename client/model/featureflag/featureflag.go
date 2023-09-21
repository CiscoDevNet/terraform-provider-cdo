package featureflag

import "strings"

type Type string

const (
	AsaConfigurationObjectMigration Type = "asa_configuration_object_migration"
)

func (t Type) String() string {
	return strings.ToLower(string(t))
}
