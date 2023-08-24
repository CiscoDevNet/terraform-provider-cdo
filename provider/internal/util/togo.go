package util

import "github.com/hashicorp/terraform-plugin-framework/types"

func TFStringListToGoStringList(l types.List) []string {
	res := make([]string, len(l.Elements()))
	for i, v := range l.Elements() {
		res[i] = v.String()
	}
	return res
}
