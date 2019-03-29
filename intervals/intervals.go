package intervals

import (
	"fmt"
	"sort"
	"time"
)

type SortedIntervals interface {
	sort.Interface

	EndBeforeStart(i int, j int) bool
}

type TimePeriod interface {
	PeriodStart() time.Time
	PeriodEnd() time.Time
}

type TimePeriods interface {
	PeriodAt(i int) TimePeriod
	Len() int
}

func CheckOverlap(si SortedIntervals) error {
	n := si.Len()
	for i := 0; i < n-1; i++ {
		// Check if intervals are sorted (on the fly)
		if !si.Less(i, i+1) {
			return fmt.Errorf("[%d].start must be after [%d].start", i+1, i)
		}

		if !si.EndBeforeStart(i, i+1) {
			return fmt.Errorf("[%d].start must be after [%d].end", i+1, i)
		}
	}
	return nil
}

func PeriodsInRange(periods TimePeriods, rng TimePeriod) (i, j int) {
	n := periods.Len()

	// Bin search from front
	i = sort.Search(n, func(x int) bool {
		return InRange(periods.PeriodAt(x), rng)
	})
	if i == n {
		return 0, 0
	}

	// Bin search from back (the remainder)
	j = n - sort.Search(n-i, func(y int) bool {
		return InRange(periods.PeriodAt(n-y-1), rng)
	})

	return i, j
}

// InRange tests if p is in rng such that the intersection of both is not empty
func InRange(p TimePeriod, rng TimePeriod) bool {
	// this check was simplified from 4 possible cases of overlap by negating the checks
	return p.PeriodStart().Before(rng.PeriodEnd()) && p.PeriodEnd().After(rng.PeriodStart())
}
