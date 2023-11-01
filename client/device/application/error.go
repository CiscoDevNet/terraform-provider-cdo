package application

import "fmt"

var (
	NotFoundError = fmt.Errorf("unable to find application")
)
