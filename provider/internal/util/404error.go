package util

import (
	"net/http"
	"strings"
)

func Is404Error(err error) bool {
	return strings.Contains(err.Error(), http.StatusText(http.StatusNotFound))
}
