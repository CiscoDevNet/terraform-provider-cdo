package util

import "strings"

func Is404Error(err error) bool {
	return strings.Contains(err.Error(), "code=404")
}
