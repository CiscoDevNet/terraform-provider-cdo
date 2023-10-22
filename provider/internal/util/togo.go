package util

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
