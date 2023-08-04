package jsonutil

import "encoding/json"

func UnmarshalStruct[T any](bytes []byte) (*T, error) {
	var out T

	if err := json.Unmarshal(bytes, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
