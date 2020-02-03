package utils

import "time"

func GetLocalTime(t time.Time) time.Time {
	return t.In(time.Local)
}

func GetTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func GetLocalTimeString(t time.Time) string {
	t = GetLocalTime(t)
	return GetTimeString(t)
}
