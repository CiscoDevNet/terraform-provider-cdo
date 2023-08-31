package statemachine

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/cdo"
)

var (
	NotFoundError           = fmt.Errorf("statemachine instance not found")
	MoreThanOneRunningError = fmt.Errorf("multiple running instances found, this is not an expected error, please report this issue at: %s", cdo.TerraformProviderCDOIssuesUrl)
)
