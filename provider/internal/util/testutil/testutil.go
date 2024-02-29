package testutil

import (
	"encoding/json"
	"fmt"
)

func CheckEqual(expected string) func(value string) error {
	return func(value string) error {
		if value != expected {
			return fmt.Errorf("string is not equal to expected value: %s (value) != %s (expected)", value, expected)
		}
		return nil
	}
}

func MustJson(input any) string {
	output, err := json.Marshal(input)
	if err != nil {
		panic("unable to marshall json for")
	}

	return string(output)
}
