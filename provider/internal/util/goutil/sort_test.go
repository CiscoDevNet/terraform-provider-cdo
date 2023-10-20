package goutil_test

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/internal/util/goutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortStrings1(t *testing.T) {
	assertSortStrings(
		t,
		[]string{"1", "2", "3"},
		[]string{"1", "2", "3"},
	)
}

func TestSortStrings2(t *testing.T) {
	assertSortStrings(
		t,
		[]string{"2", "3", "1"},
		[]string{"1", "2", "3"},
	)
}

func TestSortStrings3(t *testing.T) {
	assertSortStrings(
		t,
		[]string{"2"},
		[]string{"2"},
	)
}

func TestSortStrings4(t *testing.T) {
	assertSortStrings(
		t,
		[]string{"c", "a", "b", "c", "b", "a", "a", "b", "c"},
		[]string{"a", "a", "a", "b", "b", "b", "c", "c", "c"},
	)
}

func TestSortStrings5(t *testing.T) {
	assertSortStrings(
		t,
		[]string{"bbbz", "aac", "aaa", "bbby", "aab", "bbbb", "aad"},
		[]string{"aaa", "aab", "aac", "aad", "bbbb", "bbby", "bbbz"},
	)
}

func TestSortStrings_nilBecomesEmptySlice(t *testing.T) {
	assertSortStrings(t, nil, []string{})
}

func TestSortStrings_EmptySliceBecomesEmptySlice(t *testing.T) {
	assertSortStrings(t, []string{}, []string{})
}

func assertSortStrings(t *testing.T, unsorted []string, sorted []string) {
	assertInputNotChanged(t, unsorted)
	assertEqualAfterSort(t, unsorted, sorted)
}

func assertEqualAfterSort(t *testing.T, unsorted []string, sorted []string) {
	// call implementation
	actualSorted := goutil.SortStrings(unsorted)
	// make sure output is the same as expected sorted
	assert.Equal(t, actualSorted, sorted, fmt.Sprintf("output is not sorted, actual=%+v, expected=%+v", actualSorted, sorted))
}
func assertInputNotChanged(t *testing.T, unsorted []string) {
	// make a copy of the unsorted string
	var unsortedCopy []string
	if unsorted == nil {
		unsortedCopy = nil
	} else {
		unsortedCopy = make([]string, len(unsorted))
		copy(unsortedCopy, unsorted)
	}

	// call implementation
	goutil.SortStrings(unsorted)

	// make sure input is the not changed
	assert.Equal(t, unsortedCopy, unsorted, fmt.Sprintf("input should not be altered, before=%+v, after=%+v", unsortedCopy, unsorted))
}
