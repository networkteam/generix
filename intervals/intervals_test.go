package intervals_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/networkteam/generix/intervals"
)

// max int64 value minus number of seconds between Year 1 and 1970 (62135596800 seconds)
var maxTime = time.Unix(1<<63-62135596801, 999999999)

type timeInterval struct {
	start time.Time
	end   *time.Time
}

func (ti timeInterval) normEnd() time.Time {
	if ti.end != nil {
		return *ti.end
	}

	return maxTime
}

type timeIntervals []timeInterval

func (tis timeIntervals) PeriodAt(i int) intervals.TimePeriod {
	return tis[i]
}

func (tis timeIntervals) Len() int {
	return len(tis)
}

// start of i is less than start of j
func (tis timeIntervals) Less(i, j int) bool {
	return tis[i].start.Before(tis[j].start)
}

func (tis timeIntervals) Swap(i, j int) {
	tis[i], tis[j] = tis[j], tis[i]
}

// end of i is less than start of j
func (tis timeIntervals) EndBeforeStart(i, j int) bool {
	return tis[i].normEnd().Before(tis[j].start)
}

func (ti timeInterval) PeriodStart() time.Time {
	return ti.start
}

func (ti timeInterval) PeriodEnd() time.Time {
	return ti.normEnd()
}

var _ intervals.SortedIntervals = timeIntervals{}

func TestTimeIntervals(t *testing.T) {
	tt := []struct {
		name               string
		tis                timeIntervals
		expectedErrMessage string
	}{
		{
			name: "single open interval",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   nil,
				},
			},
			expectedErrMessage: "",
		},
		{
			name: "intervals not sorted",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
					end:   nil,
				},
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 5, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectedErrMessage: "[1].start must be after [0].start",
		},
		{
			name: "t1.end overlaps t2.start",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 8, 0, 0, 0, 0, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 5, 0, 0, 0, 0, time.UTC),
					end:   nil,
				},
			},
			expectedErrMessage: "[1].start must be after [0].end",
		},
		{
			name: "t2 inside t1",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 20, 0, 0, 0, 0, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 8, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 12, 0, 0, 0, 0, time.UTC)),
				},
			},
			expectedErrMessage: "[1].start must be after [0].end",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			var si intervals.SortedIntervals = tc.tis

			err := intervals.CheckOverlap(si)

			if tc.expectedErrMessage == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)

				assert.Equal(t, tc.expectedErrMessage, err.Error())
			}
		})
	}
}

func TestPeriodInRange(t *testing.T) {
	rng := timeInterval{
		// Beginning of January, 8th
		start: time.Date(2019, time.January, 8, 0, 0, 0, 0, time.UTC),
		// At end of January, 19th
		end: timePtr(time.Date(2019, time.January, 20, 0, 0, 0, -1, time.UTC)),
	}

	tt := []struct {
		name      string
		ti        timeInterval
		isInRange bool
	}{
		{
			name: "period outside of range",
			ti: timeInterval{
				start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 23, 0, 0, 0, -1, time.UTC)),
			},
			isInRange: true,
		},
		{
			name: "period inside in range",
			ti: timeInterval{
				start: time.Date(2019, time.January, 9, 0, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 19, 0, 0, 0, -1, time.UTC)),
			},
			isInRange: true,
		},
		{
			name: "period inside in range, start equal",
			ti: timeInterval{
				start: rng.start,
				end:   timePtr(time.Date(2019, time.January, 19, 0, 0, 0, -1, time.UTC)),
			},
			isInRange: true,
		},
		{
			name: "period inside in range, end equal",
			ti: timeInterval{
				start: time.Date(2019, time.January, 9, 0, 0, 0, 0, time.UTC),
				end:   rng.end,
			},
			isInRange: true,
		},
		{
			name:      "period equal to range",
			ti:        rng,
			isInRange: true,
		},
		{
			name: "period cuts range at start",
			ti: timeInterval{
				start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 8, 4, 0, 0, 0, time.UTC)),
			},
			isInRange: true,
		},
		{
			name: "period cuts range at end",
			ti: timeInterval{
				start: time.Date(2019, time.January, 12, 0, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 20, 0, 0, 0, 0, time.UTC)),
			},
			isInRange: true,
		},
		{
			name: "period end before range start",
			ti: timeInterval{
				start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 6, 0, 0, 0, 0, time.UTC)),
			},
			isInRange: false,
		},
		{
			name: "period end at start of range",
			ti: timeInterval{
				start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 8, 0, 0, 0, 0, time.UTC)),
			},
			isInRange: false,
		},
		{
			name: "period start equal to range end",
			ti: timeInterval{
				start: *rng.end,
				end:   timePtr(time.Date(2019, time.January, 23, 0, 0, 0, -1, time.UTC)),
			},
			isInRange: false,
		},
		{
			name: "period start after range end",
			ti: timeInterval{
				start: time.Date(2019, time.January, 20, 2, 0, 0, 0, time.UTC),
				end:   timePtr(time.Date(2019, time.January, 23, 0, 0, 0, -1, time.UTC)),
			},
			isInRange: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			isInRange := intervals.InRange(tc.ti, rng)
			assert.Equal(t, tc.isInRange, isInRange)
		})
	}
}

func TestPeriodsInRange(t *testing.T) {
	rng := timeInterval{
		// Beginning of January, 8th
		start: time.Date(2019, time.January, 8, 0, 0, 0, 0, time.UTC),
		// At end of January, 19th
		end: timePtr(time.Date(2019, time.January, 20, 0, 0, 0, -1, time.UTC)),
	}

	tt := []struct {
		name string
		tis  timeIntervals
		i, j int
	}{
		{
			name: "empty",
			tis:  timeIntervals{},
			i:    0,
			j:    0,
		},
		{
			name: "{}",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.February, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.February, 7, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.February, 7, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.February, 21, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.February, 22, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.February, 23, 0, 0, 0, -1, time.UTC)),
				},
			},
			i: 0,
			j: 0,
		},
		{
			name: "{1}",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 7, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 7, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 21, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 22, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 23, 0, 0, 0, -1, time.UTC)),
				},
			},
			i: 1,
			j: 2,
		},
		{
			name: "{1,2}",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 7, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 7, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 19, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 19, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 20, 0, 0, 0, -1, time.UTC)),
				},
			},
			i: 1,
			j: 3,
		},
		{
			name: "{0,1,2}",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 9, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 9, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 19, 0, 0, 0, -1, time.UTC)),
				},
				{
					start: time.Date(2019, time.January, 19, 0, 0, 0, 0, time.UTC),
					end:   timePtr(time.Date(2019, time.January, 20, 0, 0, 0, -1, time.UTC)),
				},
			},
			i: 0,
			j: 3,
		},
		{
			name: "{0}",
			tis: timeIntervals{
				{
					start: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
					end:   nil,
				},
			},
			i: 0,
			j: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, intervals.CheckOverlap(tc.tis), "intervals must not overlap")

			i, j := intervals.PeriodsInRange(tc.tis, rng)
			assert.Equal(t, tc.i, i, "i")
			assert.Equal(t, tc.j, j, "j")
		})
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}
