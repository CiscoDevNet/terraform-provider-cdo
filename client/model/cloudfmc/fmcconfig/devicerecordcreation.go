package fmcconfig

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
)

type DeviceRecordCreation struct {
	//Links internal.Links             `json:"links"`
	//Items []DeviceRecordCreationItem `json:"items"`
	Type            string `json:"type"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	RegKey          string `json:"regKey"`
	PerformanceTier string `json:"performanceTier"`
	NatID           string `json:"natID"`
	KeepLocalEvents bool   `json:"keepLocalEvents"`

	AccessPolicy accessPolicy    `json:"accessPolicy"`
	LicenseCaps  *[]license.Type `json:"license_caps"`
	Metadata     metadata        `json:"metadata"`
}

type accessPolicy struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type metadata struct {
	Task              task `json:"task"`
	IsPartOfContainer bool `json:"isPartOfContainer"`
	IsMultiInstance   bool `json:"isMultiInstance"`
}

type task struct {
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

//{
//    "type": "Device",
//    "name": "test-wl-ftd",
//    "version": "6.0.1",
//    "accessPolicy": {
//        "id": "06AE8B8C-5F91-0ed3-0000-004294967346",
//        "type": "AccessPolicy"
//    },
//    "regKey": "NVSr9JszUQh69oNTHjRk8ztXdtzkwurX",
//    "license_caps": [
//        "BASE"
//    ],
//    "performanceTier": "FTDv5",
//    "natID": "Vbufy7GTmZK0xXvzs3WEHZC2ueXChwUH",
//    "keepLocalEvents": false,
//    "metadata": {
//        "task": {
//            "name": "DEVICE_REGISTRATION",
//            "id": "4294972234",
//            "type": "TaskStatus"
//        },
//        "isPartOfContainer": false,
//        "isMultiInstance": false
//    },
