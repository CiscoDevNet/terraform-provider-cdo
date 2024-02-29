package testing

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/device/tags"
)

func NewTestingTags() tags.Type {
	return tags.New(
		[]string{"label-1", "label-2", "label-3"},
		map[string][]string{"grouped-labels": {"grouped-label-1", "grouped-label-2", "grouped-label-3"}},
	)
}
