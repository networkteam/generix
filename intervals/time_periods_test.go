package intervals_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/networkteam/generix/intervals"
)

type timePeriod struct {
	start time.Time
	end   time.Time
}

func (p timePeriod) String() string {
	return fmt.Sprintf("[%s, %s]", p.start, p.end)
}

func (p timePeriod) PeriodStart() time.Time {
	return p.start
}

func (p timePeriod) PeriodEnd() time.Time {
	return p.end
}

func TestPeriodDuration(t *testing.T) {
	tests := []struct {
		name string
		p    timePeriod
		want time.Duration
	}{
		{
			name: "end after start",
			p: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
			want: mustDuration(time.ParseDuration("13h55m12s")),
		},
		{
			name: "start after end",
			p: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 18, 00, 04, 42, 0, time.UTC),
			},
			want: mustDuration(time.ParseDuration("-277h43m18s")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := intervals.PeriodDuration(tt.p); got != tt.want {
				t.Errorf("PeriodDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncatePeriod(t *testing.T) {
	tests := []struct {
		name string
		p    timePeriod
		rng  timePeriod
		want timePeriod
	}{
		{
			name: "fully in range",
			p: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
			rng: timePeriod{
				start: time.Date(2019, time.March, 01, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 31, 0, 0, 0, -1, time.UTC),
			},
			want: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
		},
		{
			name: "end outside range",
			p: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
			rng: timePeriod{
				start: time.Date(2019, time.March, 29, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 0, 0, 0, -1, time.UTC),
			},
			want: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 0, 0, 0, -1, time.UTC),
			},
		},
		{
			name: "start outside range",
			p: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
			rng: timePeriod{
				start: time.Date(2019, time.March, 30, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 31, 0, 0, 0, -1, time.UTC),
			},
			want: timePeriod{
				start: time.Date(2019, time.March, 30, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
		},
		{
			name: "completely out of range",
			p: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
			rng: timePeriod{
				start: time.Date(2019, time.January, 0, 0, 0, 0, 0, time.UTC),
				end:   time.Date(2019, time.January, 31, 0, 0, 0, -1, time.UTC),
			},
			want: timePeriod{
				start: time.Date(2019, time.March, 29, 13, 48, 0, 0, time.UTC),
				end:   time.Date(2019, time.March, 30, 03, 43, 12, 0, time.UTC),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := intervals.TruncatePeriod(tt.p, tt.rng)
			actual := timePeriod{start, end}
			if actual != tt.want {
				t.Errorf("TruncatePeriod() = %v, want %v", actual, tt.want)
			}
		})
	}
}

func mustDuration(duration time.Duration, err error) time.Duration {
	if err != nil {
		panic(err)
	}
	return duration
}
