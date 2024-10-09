package utils

import (
	"errors"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/trace"
	"regexp"
	"strconv"
	"time"
)

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

func GetTimeUnitParts(timeUnit string) (num int, unit string, err error) {
	re := regexp.MustCompile(`^(\d+)([a-zA-Z])$`)
	groups := re.FindStringSubmatch(timeUnit)
	if len(groups) < 3 {
		err = errors.New("failed to parse duration text")
		log.Errorf("failed to parse duration text: %v", err)
		trace.PrintError(err)
		return 0, "", err
	}
	num, err = strconv.Atoi(groups[1])
	if err != nil {
		log.Errorf("failed to convert string to int: %v", err)
		trace.PrintError(err)
		return 0, "", err
	}
	unit = groups[2]
	return num, unit, nil
}

func GetTimeDuration(num string, unit string) (d time.Duration, err error) {
	numInt, err := strconv.Atoi(num)
	if err != nil {
		log.Errorf("failed to convert string to int: %v", err)
		trace.PrintError(err)
		return d, err
	}
	switch unit {
	case "s":
		d = time.Duration(numInt) * time.Second
	case "m":
		d = time.Duration(numInt) * time.Minute
	case "h":
		d = time.Duration(numInt) * time.Hour
	case "d":
		d = time.Duration(numInt) * 24 * time.Hour
	case "w":
		d = time.Duration(numInt) * 7 * 24 * time.Hour
	case "M":
		d = time.Duration(numInt) * 30 * 24 * time.Hour
	case "y":
		d = time.Duration(numInt) * 365 * 24 * time.Hour
	default:
		err = errors.New("invalid time unit")
		log.Errorf("invalid time unit: %v", unit)
		trace.PrintError(err)
		return d, err
	}
	return d, nil
}
