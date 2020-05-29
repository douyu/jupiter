package xtime

import "time"

// GetTimestampInMilli ...
func GetTimestampInMilli() int64 {
	return int64(time.Now().UnixNano() / 1e6)
}

// Timing the cost of function call, unix nano was returned
func Elapse(f func()) int64 {
	now := time.Now().UnixNano()
	f()
	return time.Now().UnixNano() - now
}

// IsLeapYear ...
func IsLeapYear(year int) bool {
	if year%100 == 0 {
		return year%400 == 0
	}

	return year%4 == 0
}
