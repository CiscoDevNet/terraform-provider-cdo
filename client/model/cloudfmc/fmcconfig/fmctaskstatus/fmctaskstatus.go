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
	if len(b) <= 2 || b == nil {
		return fmt.Errorf("cannot unmarshal empty tring as a fmc task status, it should be one of valid roles: %+v", NameToTypeMap)
	}
	type_, ok := NameToTypeMap[string(b[1:len(b)-1])] // strip off quote
	if !ok {
		return fmt.Errorf("failed to unmarshal %s into FMC task status, should be one of %+v", type_, NameToTypeMap)
	}
	*t = type_

	return nil
}
