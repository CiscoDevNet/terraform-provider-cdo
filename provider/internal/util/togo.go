package util

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TFStringListToGoStringList(l []types.String) []string {
	return sliceutil.Map(l, func(tfString types.String) string {
		return tfString.ValueString()
	})
}
