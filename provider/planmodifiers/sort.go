package planmodifiers

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/sliceutil"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func UseStateForUnorderedStringList() planmodifier.List {
	return useStateForUnorderedStringList{}
}

// useStateForUnorderedStringList
type useStateForUnorderedStringList struct{}

// Description returns a human-readable description of the plan modifier.
func (m useStateForUnorderedStringList) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m useStateForUnorderedStringList) MarkdownDescription(_ context.Context) string {
	return "If the list are the same after sorting, use the existing state value."
}

// PlanModifyList implements the plan modification logic.
func (m useStateForUnorderedStringList) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {

	// no need to compare on resource destroy.
	if req.Plan.Raw.IsNull() {
		return
	}

	// no need to compare on resource creation
	if req.State.Raw.IsNull() {
		return
	}

	// only compare on resource update.

	// read plan value
	planElements := make([]types.String, 0, len(req.PlanValue.Elements()))
	resp.Diagnostics.Append(req.PlanValue.ElementsAs(ctx, &planElements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	planList := util.TFStringListToGoStringList(planElements)

	// read state value
	stateElements := make([]types.String, 0, len(req.StateValue.Elements()))
	resp.Diagnostics.Append(req.PlanValue.ElementsAs(ctx, &stateElements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}
	stateList := util.TFStringListToGoStringList(stateElements)

	// compare and ignore plan if they are equal, just maybe in unordered value
	if sliceutil.StringsEqualUnordered(planList, stateList) {
		// discard changes in plan by overwriting it with state value
		resp.PlanValue = req.StateValue
	}
}
