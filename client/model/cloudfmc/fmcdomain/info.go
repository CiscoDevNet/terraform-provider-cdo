package fmcdomain

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/common"
)

type Info struct {
	Links  common.Links  `json:"links"`
	Paging common.Paging `json:"paging"`
	Items  []Item        `json:"items"`
}

type Item struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
}
