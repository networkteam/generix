package intervals

import (
	"fmt"
	"sort"
)

// SortableIntervals represents sortable intervals that have a start and end value that can be compared
type SortableIntervals interface {
	// Implements Less by comparing the start values
	sort.Interface

	// EndBeforeStart checks if the end of interval at index i is strictly before
	// the start of interval at index j, where j>i
	EndBeforeStart(i int, j int) bool
}

// CheckOverlap checks if any two intervals are overlapping
func CheckOverlap(si SortableIntervals) error {
	n := si.Len()
	for i := 0; i < n-1; i++ {
		// Make sure intervals are sorted (on the fly)
		if !si.Less(i, i+1) {
			return fmt.Errorf("[%d].start must be after [%d].start", i+1, i)
		}

		if !si.EndBeforeStart(i, i+1) {
			return fmt.Errorf("[%d].start must be after [%d].end", i+1, i)
		}
	}
	return nil
}
