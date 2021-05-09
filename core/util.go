package core

import (
	"github.com/pkg/errors"
	"time"
)

func ParseDateWithTime(dateYYYYMMDD, timeHHMMSS string) (time.Time, error) {
	parsee := dateYYYYMMDD + "T" + timeHHMMSS + "Z"
	parsed, err := time.Parse(time.RFC3339, parsee)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "could not parse time '%s'", parsee)
	}

	return parsed, nil
}
