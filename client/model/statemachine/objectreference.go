package statemachine

type ObjectReference struct {
	Namespace string `json:"namespace"`
	Type      string `json:"type"`
	Uid       string `json:"uid"`
}
