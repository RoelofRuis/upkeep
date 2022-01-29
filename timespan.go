package main

import "time"

type TimeBlock struct {
	Start Moment
	End   Moment
	Tags  []string
}

func NewTimeBlock(t time.Time) TimeBlock {
	return TimeBlock{
		Start: NewMoment(t),
		End:   Moment{},
		Tags:  []string{},
	}
}

func (ts *TimeBlock) AddTag(t string) {
	for _, tag := range ts.Tags {
		if tag == t {
			return
		}
	}
	ts.Tags = append(ts.Tags, t)
}

func (ts *TimeBlock) RemoveTag(t string) {
	for i, tag := range ts.Tags {
		if tag == t {
			ts.Tags[i] = ts.Tags[len(ts.Tags)-1]
			ts.Tags = ts.Tags[:len(ts.Tags)-1]
			return
		}
	}
}

func (ts *TimeBlock) Complete(t time.Time) {
	ts.End = NewMoment(t)
}

func (ts TimeBlock) HasEnded() bool {
	return ts.End.IsDefined()
}
