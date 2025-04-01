package util

import "time"

func FromBoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func FromIntToBool(value int) bool {
	return value != 0
}

func FromFloat64ToTime(value float64) time.Time {
	seconds := int64(value)
	nanoseconds := int64((value - float64(value)) * 1e9)
	return time.Unix(seconds, nanoseconds)
}

func FromTimeToFloat64(value time.Time) float64 {
	return float64(value.Unix())
}
