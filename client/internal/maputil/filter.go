package maputil

// FilterKeys creates a copy of the passed map containing only keys that return true when passed into the predicate
func FilterKeys[Key comparable, Value any](m map[Key]Value, predicate func(Key) bool) map[Key]Value {
	resultMap := map[Key]Value{}

	for k, v := range m {
		if predicate(k) {
			resultMap[k] = v
		}
	}

	return resultMap
}
