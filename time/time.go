package time

import (
	"encoding/base64"
	"time"
)

var timeFormat = "2006-01-02T15:04:05.999Z07:00"

func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)
	return base64.StdEncoding.EncodeToString([]byte(timeString))
}

func DecodeCursor(encodedTime string) (time.Time, error) {
	b, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}
	timeString := string(b)
	return time.Parse(timeFormat, timeString)
}
