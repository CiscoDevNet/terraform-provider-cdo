package fmcconfig

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"

type DeviceRecordCreation struct {
	Links internal.Links             `json:"links"`
	Items []DeviceRecordCreationItem `json:"items"`
}

type DeviceRecordCreationItem struct {
	Id              string         `json:"id"`
	Type            string         `json:"type"`
	Links           internal.Links `json:"links"`
	Name            string         `json:"name"`
	Version         string         `json:"version"`
	KeepLocalEvents bool           `json:"keepLocalEvents"`
}
