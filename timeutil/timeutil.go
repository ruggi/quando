package timeutil

import "time"

var Now = func() time.Time {
	return time.Now()
}

func Today() time.Time {
	now := Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
