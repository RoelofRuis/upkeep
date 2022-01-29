package main

import "time"

type TimeBlock struct {
	Start Moment
	End   Moment
}

func NewTimeBlock(t time.Time) TimeBlock {
	return TimeBlock{
		Start: NewMoment(t),
		End:   Moment{},
	}
}

func (ts *TimeBlock) Complete(t time.Time) {
	ts.End = NewMoment(t)
}

func (ts TimeBlock) HasEnded() bool {
	return ts.End.IsDefined()
}
