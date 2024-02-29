package util

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func TFStringListToGoStringList(l []types.String) []string {
	return sliceutil.Map(l, func(tfString types.String) string {
		return tfString.ValueString()
	})
}

func TFStringSetToGoStringList(ctx context.Context, s types.Set) ([]string, error) {
	elements := s.Elements()
	stringList := make([]string, len(elements))
	d := s.ElementsAs(ctx, &stringList, false)
	if d.HasError() {
		return nil, fmt.Errorf(DiagSummary(d))
	}
	return stringList, nil
}

func ToLabels(ctx context.Context, labels types.Set, groupedLabels types.Map) (tags.Type, error) {
	convertedLabels, err := TFStringSetToGoStringList(ctx, labels)
	if err != nil {
		return nil, err
	}

	convertedGroupLabels, err := TFMapToGoMapOfStringSlices(ctx, groupedLabels)
	if err != nil {
		return nil, err
	}

	return tags.New(convertedLabels, convertedGroupLabels), nil
}

func TFStringSetToLicenses(ctx context.Context, s types.Set) ([]license.Type, error) {
	licensesGoList, err := TFStringSetToGoStringList(ctx, s)
	if err != nil {
		return nil, err
	}
	licenses, err := license.StringToCdoLicenses(strings.Join(licensesGoList, ","))
	if err != nil {
		return nil, err
	}
	return licenses, nil
}

func TFMapToGoMap[GoType any](m types.Map, elemConverter func(attr.Value) (GoType, error)) (map[string]GoType, error) {
	resultMap := map[string]GoType{}

	for k, v := range m.Elements() {
		elem, err := elemConverter(v)
		if err != nil {
			return nil, err
		}

		resultMap[k] = elem
	}

	return resultMap, nil
}

func TFMapToGoMapOfStringSlices(ctx context.Context, m types.Map) (map[string][]string, error) {
	return TFMapToGoMap[[]string](m, func(v attr.Value) ([]string, error) {
		n, ok := v.(basetypes.SetValue)
		if !ok {
			return nil, errors.New("unexpected element type in tf map value")
		}

		return TFStringSetToGoStringList(ctx, n)
	})
}
