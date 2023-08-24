package sliceutil

// MapWithError takes a slice and a function, apply the function to each element in the slice, and return the result.
// It takes a function that can return error, when it returns error, no subsequent element is mapped
// and the slice mapped so far is returned together with the error.
func MapWithError[T any, V any](sliceT []T, mapFunc func(T) (V, error)) ([]V, error) {
	sliceV := make([]V, len(sliceT))
	for i, item := range sliceT {
		mapped, err := mapFunc(item)
		if err != nil {
			return sliceV, err
		}
		sliceV[i] = mapped
	}
	return sliceV, nil
}

// Map takes a slice and a function, apply the function to each element in the slice, and return the mapped slice.
func Map[T any, V any](sliceT []T, mapFunc func(T) V) []V {
	sliceV := make([]V, len(sliceT))
	for i, item := range sliceT {
		sliceV[i] = mapFunc(item)
	}
	return sliceV
}
