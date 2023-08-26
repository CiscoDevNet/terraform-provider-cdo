package fmcdomain

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc"

type Info struct {
	Links  cloudfmc.Links  `json:"links"`
	Paging cloudfmc.Paging `json:"paging"`
	Items  []Item          `json:"items"`
}

type Item struct {
	Uuid string `json:"uuid"`
	Name string `json:"name"`
	Type string `json:"type"`
}
