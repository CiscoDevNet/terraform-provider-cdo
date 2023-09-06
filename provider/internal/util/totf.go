package util

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GoStringSliceToTFStringList(stringSlice []string) []types.String {
	return sliceutil.Map(stringSlice, func(s string) types.String {
		return types.StringValue(s)
	})
}
