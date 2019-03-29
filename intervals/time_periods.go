package intervals

import (
	"sort"
	"time"
)

// TimePeriod has a start and end time.Time value, where end should be after start
type TimePeriod interface {
	PeriodStart() time.Time
	PeriodEnd() time.Time
}

// TimePeriods is a list of TimePeriod
type TimePeriods interface {
	// Get period at index i
	PeriodAt(i int) TimePeriod
	Len() int
}

// PeriodDuration gives the time.Duration of a period
func PeriodDuration(p TimePeriod) time.Duration {
	return p.PeriodEnd().Sub(p.PeriodStart())
}

// TruncatePeriod returns start end end times of a period that are inside the given range
//
// Start and end of period are returned if period is not in the given range.
func TruncatePeriod(p TimePeriod, rng TimePeriod) (start time.Time, end time.Time) {
	start, end = p.PeriodStart(), p.PeriodEnd()
	if !InRange(p, rng) {
		return
	}
	if p.PeriodStart().Before(rng.PeriodStart()) {
		start = rng.PeriodStart()
	}
	if p.PeriodEnd().After(rng.PeriodEnd()) {
		end = rng.PeriodEnd()
	}
	return
}

// PeriodsOverlapRange returns start and end index (exclusive) for all
// periods overlapping the given range:
//
// - periods must be sorted
//
// - periods itself must not overlap
//
// Note: A period lying between the given indices could start before the range
// or end after the range.
func PeriodsOverlapRange(periods TimePeriods, rng TimePeriod) (i, j int) {
	n := periods.Len()

	// Bin search from front
	i = sort.Search(n, func(x int) bool {
		return InRange(periods.PeriodAt(x), rng)
	})
	// No element found (invariant of sort.Search)
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
	// this check is deduced by transforming the negation of InRange (all periods not in range)
	return p.PeriodStart().Before(rng.PeriodEnd()) && p.PeriodEnd().After(rng.PeriodStart())
}
