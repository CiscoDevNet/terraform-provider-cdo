package fmctaskstatus

import "fmt"

type Type string

const (
	Running Type = "RUNNING"
	Success Type = "SUCCESS"
	Failed  Type = "FAILED"
	Pending Type = "PENDING"
)

var NameToTypeMap = map[string]Type{
	"RUNNING": Running,
	"SUCCESS": Success,
	"FAILED":  Failed,
	"PENDING": Pending,
}

func (t *Type) MarshalJSON() ([]byte, error) {
	return []byte(*t), nil
}

func (t *Type) UnmarshalJSON(b []byte) error {
	type_, ok := NameToTypeMap[string(b)]
	if !ok {
		return fmt.Errorf("failed to unmarshal %s into FMC task status, should be one of %+v", type_, NameToTypeMap)
	}
	*t = type_

	return nil
}
