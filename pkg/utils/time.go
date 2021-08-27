package utils

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type relTimeMagnitude struct {
	D     time.Duration
	Name  string
	DivBy time.Duration
	ModBy time.Duration
}

const (
	Day   = time.Hour * 24
	Week  = Day * 7
	Month = Day * 30
	Year  = Month * 12
)

var magnitudes = []relTimeMagnitude{
	{time.Minute, "second", time.Second, 60},
	{time.Hour, "minute", time.Minute, 60},
	{Day, "hour", time.Hour, 24},
	{Week, "day", Day, 7},
	{Month, "week", Week, 7},
	{Year, "month", Month, 12},
	{math.MaxInt64, "year", Year, -1},
}

func CustomDurationString(diff time.Duration, numParts int, glue string) string {
	if diff < time.Second {
		return "now"
	}

	n := sort.Search(len(magnitudes), func(i int) bool {
		return magnitudes[i].D > diff
	})

	if n >= len(magnitudes) {
		n--
	}

	var parts []string

	partIndex := 0
	for i := 0; partIndex < numParts && n-i >= 0; i++ {
		mag := magnitudes[n-i]

		value := diff
		if mag.DivBy != -1 {
			value /= mag.DivBy
		}
		if mag.ModBy != -1 {
			value %= mag.ModBy
		}
		if value > 0 {
			part := fmt.Sprintf("%d %s", value, mag.Name)
			if value > 1 {
				part += "s"
			}

			diff -= value * mag.DivBy

			parts = append(parts, part)
			partIndex++
		}
	}

	return strings.Join(parts, glue)
}

func TimeSince(t2 time.Time) string {
	t1 := time.Now()
	//return CustomRelTime(t1, t2, 2, " ")
	return CustomDurationString(t1.Sub(t2), 3, " ")
}
