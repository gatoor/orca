package util

import "time"

func GetNow() (string, time.Time) {
	t := time.Now().UTC()
	return t.Format(time.RFC3339Nano), t
}