package sliceutil

import "sort"

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

// SortStrings is a non-in-place version of sort.Strings
func SortStrings(toSort []string) []string {
	temp := make([]string, len(toSort))
	copy(temp, toSort)
	sort.Strings(temp)
	return temp
}
