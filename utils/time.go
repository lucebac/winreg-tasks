// SPDX-License-Identifier: MIT

package utils

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/lucebac/winreg-tasks/generated"
)

const secondsUntilEpoch = 11_644_473_600

var dateMin = time.Date(1601, 1, 1, 0, 0, 0, 0, time.UTC)
var dateMax = time.Date(9999, 12, 31, 23, 59, 59, 999999999, time.UTC)

func TimeFromFILETIME(filetime uint64) time.Time {
	epoch := int64(filetime/10_000_000) - secondsUntilEpoch

	// we need to cap negative values (any datetime before 1970-01-01)
	// and excessively large positive values (any datetime after 9999-12-31T23:59:59.99999999)
	// because otherwise JSON serialization of timestamps fails
	if epoch < 0 {
		return dateMin
	} else if epoch >= dateMax.Unix() {
		return dateMax
	}

	return time.Unix(epoch, 0)
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
	if d.Duration.Seconds() >= 0xffffffff {
		return json.Marshal("infinity")
	}

	return json.Marshal(d.String())
}

func SecondsToDuration(nsec uint32) Duration {
	return Duration{Duration: time.Duration(uint64(nsec) * uint64(time.Second))}
}

const RFC3339NoTimezone = "2006-01-02T15:04:05.999999999"

func ParseWindowsTimestamp(s string) (*time.Time, error) {
	t, err := time.Parse(RFC3339NoTimezone, s)
	if err == nil {
		return &t, nil
	}

	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return &t, nil
	}

	return nil, errors.New("unexpected timestamp format")
}
