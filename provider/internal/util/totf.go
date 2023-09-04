package util

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GoStringSliceToTFStringList(stringSlice []string) []types.String {
	l := make([]types.String, len(stringSlice))
	for i, v := range stringSlice {
		l[i] = types.StringValue(v)
	}

	return l
}
