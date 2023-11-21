package util

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
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

func TFStringSetToTagLabels(ctx context.Context, s types.Set) (tags.Type, error) {
	stringList, err := TFStringSetToGoStringList(ctx, s)
	if err != nil {
		return tags.Type{}, err
	}
	return tags.New(stringList...), nil
}

func TFStringSetToLicenses(ctx context.Context, s types.Set) ([]license.Type, error) {
	licensesGoList, err := TFStringSetToGoStringList(ctx, s)
	if err != nil {
		return nil, err
	}
	licenses, err := license.DeserializeAllFromCdo(strings.Join(licensesGoList, ","))
	if err != nil {
		return nil, err
	}
	return licenses, nil
}
