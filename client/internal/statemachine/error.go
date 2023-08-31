package statemachine

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
)

var (
	StateMachineNotFoundError           = fmt.Errorf("statemachine instance not found")
	MoreThanOneStateMachineRunningError = fmt.Errorf("multiple running instance found, this is an expected error, please report this issue at: %s", cdo.TerraformProviderCDOIssuesUrl)
)
