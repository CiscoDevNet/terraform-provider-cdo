package model

import (
	"fmt"
	"strings"
)

type ConnectorType string

const (
	ConnectorTypeCloudDeviceGateway    ConnectorType = "CDG"
	ConnectorTypeSecureDeviceConnector ConnectorType = "SDC"
)

func ParseConnectorType(input string) (ConnectorType, error) {
	switch strings.ToUpper(input) {
	case string(ConnectorTypeCloudDeviceGateway):
		return ConnectorTypeCloudDeviceGateway, nil

	case string(ConnectorTypeSecureDeviceConnector):
		return ConnectorTypeSecureDeviceConnector, nil

	default:
		return "", fmt.Errorf("'%s' is not a valid connector type", input)

	}
}

func MustParseConnectorType(input string) ConnectorType {
	larType, err := ParseConnectorType(input)
	if err != nil {
		panic(err)
	}

	return larType
}
