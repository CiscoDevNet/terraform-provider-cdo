package util

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func GoStringSliceToTFStringList(stringSlice []string) []types.String {
	return sliceutil.Map(stringSlice, func(s string) types.String {
		return types.StringValue(s)
	})
}

func GoStringSliceToTFStringSet(stringSlice []string) types.Set {
	elements := sliceutil.Map(stringSlice, func(s string) attr.Value {
		return types.StringValue(s)
	})
	return types.SetValueMust(types.StringType, elements)
}

func GoMapToTFMap[GoType any](m map[string]GoType, elementType attr.Type, converter func(GoType) attr.Value) types.Map {
	tfMap := map[string]attr.Value{}

	for k, v := range m {
		tfMap[k] = converter(v)
	}

	return types.MapValueMust(elementType, tfMap)
}

func GoMapToStringSetTFMap(m map[string][]string) types.Map {
	return GoMapToTFMap(m, basetypes.SetType{ElemType: types.StringType}, func(slice []string) attr.Value {
		list := GoStringSliceToTFStringList(slice)

		listOfValues := make([]attr.Value, len(list))
		for i, v := range list {
			listOfValues[i] = v
		}

		return types.SetValueMust(types.StringType, listOfValues)
	})
}
