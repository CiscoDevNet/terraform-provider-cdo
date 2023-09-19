package acctest

import (
	"fmt"
	"os"
)

type env struct{}

var Env = &env{}

func (e *env) UserDataSourceUser() string {
	return e.mustGet("USER_DATA_SOURCE_USERNAME")
}

func (e *env) UserDataSourceRole() string {
	return e.mustGet("USER_DATA_SOURCE_ROLE")
}

func (e *env) UserDataSourceIsApiOnly() string {
	return e.mustGet("USER_DATA_SOURCE_IS_API_ONLY")
}

func (e *env) mustGet(envName string) string {
	value, ok := os.LookupEnv(envName)
	if ok {
		return value
	}
	panic(fmt.Sprintf("Acceptance test requires environment variable: %s to be set", envName))
}
