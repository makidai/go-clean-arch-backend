package repository

import (
	"encoding/base64"
	"time"
)

func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(time.RFC3339, timeString)

	return t, err
}

func EncodeCursor(t time.Time) string {
	timeString := t.Format(time.RFC3339)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
