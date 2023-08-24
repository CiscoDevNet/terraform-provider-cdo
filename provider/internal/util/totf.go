package util

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GoStringSliceToTFStringList(stringSlice []string) types.List {
	s := make([]attr.Value, len(stringSlice))
	for i, v := range stringSlice {
		s[i] = types.StringValue(v)
	}

	l, _ := types.ListValue(types.StringType, s) // drop diag as type is always correct
	return l
}
