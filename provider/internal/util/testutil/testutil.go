package testutil

import "fmt"

func CheckEqual(expected string) func(value string) error {
	return func(value string) error {
		if value != expected {
			return fmt.Errorf("string is not equal to expected value: %s (value) != %s (expected)", value, expected)
		}
		return nil
	}
}
