package goutil

import "sort"

// SortStrings return a sorted copy of the input string array.
func SortStrings(inp []string) []string {
	temp := make([]string, len(inp))
	copy(temp, inp)
	sort.Strings(temp)
	return temp
}
