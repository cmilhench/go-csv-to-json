package time

import (
	"errors"
	"time"
)

var TimeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05 -0700 MST",
	"2006-01-02T15:04:05Z",
	"2006-01-02 15:04:05.999999999 -0700 MST",
	"2006-01-02",
	"15:04:05",
}

func Parse(input string, layouts ...string) (time.Time, error) {
	var lastErr error
	for _, format := range layouts {
		if t, err := time.Parse(format, input); err == nil {
			return t, nil
		} else {
			lastErr = err
		}
	}
	return time.Time{}, errors.New("unable to parse time: " + input + " - " + lastErr.Error())
}
