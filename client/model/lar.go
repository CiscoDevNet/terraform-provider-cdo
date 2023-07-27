package model

import (
	"fmt"
	"strings"
)

type LarType string

const (
	LarTypeCloudDeviceGateway    LarType = "CDG"
	LarTypeSecureDeviceConnector LarType = "SDC"
)

func ParseLarType(input string) (LarType, error) {
	switch strings.ToUpper(input) {
	case string(LarTypeCloudDeviceGateway):
		return LarTypeCloudDeviceGateway, nil

	case string(LarTypeSecureDeviceConnector):
		return LarTypeSecureDeviceConnector, nil

	default:
		return "", fmt.Errorf("'%s' is not a valid lar type", input)

	}
}

func MustParseLarType(input string) LarType {
	larType, err := ParseLarType(input)
	if err != nil {
		panic(err)
	}

	return larType
}
