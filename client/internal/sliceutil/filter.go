package sliceutil

func Filter[T any](arr []T, function func(T) bool) []T {
	newArr := make([]T, 0)
	for _, item := range arr {
		if function(item) {
			newArr = append(newArr, item)
		}
	}
	return newArr
}
