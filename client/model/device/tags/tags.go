package tags

import (
	"encoding/json"
	"fmt"
)

// Type should be used with json:"tags"
type Type struct {
	Labels []string `json:"labels"`
}

func Empty() Type {
	return New()
}

func New(tags ...string) Type {
	return Type{
		Labels: tags,
	}
}

func (tags Type) GetLabelsJsonArrayString() string {
	b, _ := json.Marshal(tags.Labels)
	return fmt.Sprintf("%v", string(b))
}
