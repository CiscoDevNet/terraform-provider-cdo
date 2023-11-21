package goutil

import (
	"golang.org/x/exp/constraints"
	"sort"
)

// AsPointer convert interface{} to *interface{}, if input is not nil
func AsPointer(obj interface{}) *interface{} {
	var ptr *interface{} = nil
	if obj != nil {
		ptr = &obj
	}
	return ptr
}

// NewBoolPointer return a pointer of the given boolean value, this function is needed because you cant do &true or &false in golang
func NewBoolPointer(value bool) *bool {
	return &value
}

// Min return the smaller value of a and b
// starting go 1.21, we can use builtin min
func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max return the greater value of a and b
// starting go 1.21, we can use builtin max
func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// SortStrings return a sorted copy of the input string array.
func SortStrings(inp []string) []string {
	temp := make([]string, len(inp))
	copy(temp, inp)
	sort.Strings(temp)
	return temp
}
