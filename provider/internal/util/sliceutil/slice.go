package sliceutil

// Map takes a slice and a function, apply the function to each element in the slice, and return the mapped slice.
func Map[T any, V any](sliceT []T, mapFunc func(T) V) []V {
	sliceV := make([]V, len(sliceT))
	for i, item := range sliceT {
		sliceV[i] = mapFunc(item)
	}
	return sliceV
}
