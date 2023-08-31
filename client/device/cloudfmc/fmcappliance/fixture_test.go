package fmcappliance_test

const (
	baseUrl           = "https://unit-test.net"
	FmcSpecificUid    = "unit-test-cloudfmc-specific-uid"
	queueTriggerState = "unit-test-queue-trigger-state"
	uid               = "unit-test-uid"
	state             = "unit-test-state"
	domainUid         = "unit-test-domainUid"
)

var (
	stateMachineContext = map[string]string{
		"unit-test-sm-context-key": "unit-test-sm-context-value",
	}
)
