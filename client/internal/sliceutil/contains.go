package sliceutil

func Contains[T comparable](slice []T, value T) bool {
	for _, element := range slice {
		if value == element {
			return true
		}
	}

	return false
}
