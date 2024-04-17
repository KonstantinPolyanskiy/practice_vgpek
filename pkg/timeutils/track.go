package timeutils

import "time"

// TrackTime позволяет узнать, сколько времени заняло выполнение определенного кода
func TrackTime(start time.Time) time.Duration {
	return time.Since(start)
}
