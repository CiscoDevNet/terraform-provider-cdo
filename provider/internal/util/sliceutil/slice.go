package sliceutil

import (
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/goutil"
)

// Map takes a slice and a function, apply the function to each element in the slice, and return the mapped slice.
func Map[T any, V any](sliceT []T, mapFunc func(T) V) []V {
	sliceV := make([]V, len(sliceT))
	for i, item := range sliceT {
		sliceV[i] = mapFunc(item)
	}
	return sliceV
}

// MapWithError takes a slice and a function, apply the function to each element in the slice, and return the mapped slice.
// It allows the mapping function to return error, when it happens, it will terminate and return early.
func MapWithError[T any, V any](sliceT []T, mapFunc func(T) (V, error)) ([]V, error) {
	sliceV := make([]V, len(sliceT))
	for i, item := range sliceT {
		v, err := mapFunc(item)
		if err != nil {
			return nil, err
		}
		sliceV[i] = v

	}
	return sliceV, nil
}

func StringsEqual(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

func StringsEqualUnordered(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	s1Copy := goutil.SortStrings(s1)
	s2Copy := goutil.SortStrings(s2)
	return StringsEqual(s1Copy, s2Copy)
}

func Reverse[T any](s []T) []T {
	reversed := make([]T, len(s))
	if len(s) == 0 {
		return reversed
	}
	for i, j := len(s)-1, 0; i >= 0; {
		reversed[i] = s[j]
		i--
		j++
	}
	return reversed
}
