// SPDX-License-Identifier: MIT

package utils

import (
	"encoding/json"
	"time"

	"github.com/lucebac/winreg-tasks/generated"
)

const secondsUntilEpoch = 11_644_473_600

var dateMin = time.Unix(0, 0)
var dateMax = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

func TimeFromFILETIME(filetime uint64) time.Time {
	// we need to handle the two special values 0 and -1 differently to not break golang's time struct
	if filetime == 0 {
		return dateMin
	} else if filetime >= 1<<63-1 {
		return dateMax
	}

	epoch := filetime/10_000_000 - secondsUntilEpoch
	return time.Unix(int64(epoch), 0)
}

// TimeFromTSTime turns a generated TSTime object into a Golang time.Time.
// tz must be set to the Timezone of the TSTime object, otherwise the
// de-localization returns objects in the wrong timezone.
func TimeFromTSTime(gen *generated.Tstime, tz *time.Location) time.Time {
	tt := TimeFromFILETIME(uint64(gen.Filetime.HighDateTime)<<32 | uint64(gen.Filetime.LowDateTime))

	if gen.IsLocalized != 0 {
		// find out the offset to UTC by assuming tt was already in UTC and
		// translating it into the original timezone
		difference := tt.Sub(tt.In(tz))

		// now subtract the time difference so that we get the real UTC timestamp
		tt = tt.Add(-difference)
	}

	return tt
}

func DurationFromTSTimePeriod(gen *generated.Tstimeperiod) Duration {
	return Duration{Duration: time.Duration(uint64(gen.Year)*365*24*uint64(time.Hour) + // TODO: check whether Microsoft really handles a year as 365 days internally
		uint64(gen.Month)*30*24*uint64(time.Hour) + // TODO: check whether Microsoft really handles a month as 30 days internally
		uint64(gen.Day)*24*uint64(time.Hour) +
		uint64(gen.Hour)*uint64(time.Hour) +
		uint64(gen.Minute)*uint64(time.Minute) +
		uint64(gen.Second)*uint64(time.Second)),
	}
}

// time.Duration does not support JSON marshalling; hence, we need a wrapper which implements it
type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func SecondsToDuration(nsec uint32) Duration {
	return Duration{Duration: time.Duration(uint64(nsec) * uint64(time.Second))}
}
