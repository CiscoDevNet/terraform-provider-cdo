package util

import "github.com/hashicorp/terraform-plugin-framework/types"

func TFStringListToGoStringList(l []types.String) []string {
	res := make([]string, len(l))
	for i, v := range l {
		res[i] = v.ValueString()
	}
	return res
}
