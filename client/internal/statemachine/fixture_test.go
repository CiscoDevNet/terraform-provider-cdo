package statemachine_test

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/statemachine"

const (
	baseUrl   = "https://unit-test.cdo.cisco.com"
	deviceUid = "unit-test-device-uid"
)

var (
	validReadStateMachineOutput = statemachine.NewReadInstanceByDeviceUidOutputBuilder().Build()
)
