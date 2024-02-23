package publicapilabels

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/maputil"

const ungroupedLabelKeyName = "ungroupedLabels"

type Type map[string][]string

func Empty() Type {
	return Type{}
}

func NewUnlabelled(tags ...string) Type {
	return Type{
		ungroupedLabelKeyName: tags,
	}
}

func New(tags []string, groupedTags map[string][]string) Type {
	outputTags := Type{}

	for k, v := range groupedTags {
		outputTags[k] = v
	}

	outputTags[ungroupedLabelKeyName] = tags

	return outputTags
}

func (t Type) UngroupedTags() []string {
	label, ok := t[ungroupedLabelKeyName]
	if !ok {
		return []string{}
	}

	return label
}

func (t Type) GroupedTags() map[string][]string {
	return maputil.FilterKeys(t, func(s string) bool { return s != ungroupedLabelKeyName })
}
