package fmcconfig

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
)

type DeviceRecordCreation struct {
	Type            string `json:"type"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	RegKey          string `json:"regKey"`
	PerformanceTier string `json:"performanceTier"`
	NatID           string `json:"natID"`
	KeepLocalEvents bool   `json:"keepLocalEvents"`

	AccessPolicy accessPolicy    `json:"accessPolicy"`
	LicenseCaps  *[]license.Type `json:"license_caps"`
	Metadata     Metadata        `json:"metadata"`
}

type accessPolicy struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type Metadata struct {
	Task              Task `json:"task"`
	IsPartOfContainer bool `json:"isPartOfContainer"`
	IsMultiInstance   bool `json:"isMultiInstance"`
}

type Task struct {
	Name string `json:"name"`
	Id   string `json:"id"`
	Type string `json:"type"`
}

type DeviceRecordCreationItem struct {
	Id              string         `json:"id"`
	Type            string         `json:"type"`
	Links           internal.Links `json:"links"`
	Name            string         `json:"name"`
	Version         string         `json:"version"`
	KeepLocalEvents bool           `json:"keepLocalEvents"`
}
