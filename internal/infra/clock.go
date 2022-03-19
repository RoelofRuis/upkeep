package infra

import "time"

type Clock interface {
	Now() time.Time
}

type SystemClock struct{}

func (c SystemClock) Now() time.Time {
	return time.Now()
}

type FixedClock struct {
	Time time.Time
}

func (c FixedClock) Now() time.Time {
	return c.Time
}