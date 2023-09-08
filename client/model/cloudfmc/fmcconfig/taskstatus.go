package fmcconfig

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"

type TaskStatus struct {
	Id       string         `json:"id"`
	Type     string         `json:"type"`
	Links    internal.Links `json:"links"`
	TaskType string         `json:"taskType"`
	Message  string         `json:"message"`
	Status   string         `json:"status"`
	SubTasks []SubTask      `json:"subTasks"`
}

type SubTask struct {
	Id          string `json:"id"`
	Target      Target `json:"target"`
	Status      string `json:"status"`
	Message     string `json:"message"`
	Description string `json:"description"`
}

type Target struct {
	Type string `json:"type"`
	Name string `json:"name"`
}
