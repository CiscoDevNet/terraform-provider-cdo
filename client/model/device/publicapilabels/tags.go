package publicapilabels

type Type struct {
	GroupedLabels   map[string][]string `json:"groupedLabels"`
	UngroupedLabels []string            `json:"ungroupedLabels"`
}

func Empty() Type {
	return Type{}
}

func NewUnlabelled(tags ...string) Type {
	return Type{
		UngroupedLabels: tags,
		GroupedLabels:   map[string][]string{},
	}
}

func New(tags []string, groupedTags map[string][]string) Type {
	groupedLabels := map[string][]string{}

	for k, v := range groupedTags {
		groupedLabels[k] = v
	}

	return Type{
		UngroupedLabels: tags,
		GroupedLabels:   groupedLabels,
	}
}
