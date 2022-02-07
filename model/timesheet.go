package model

import (
	"time"
)

type Timesheet struct {
	Created   time.Time
	NextId    int
	Blocks    []TimeBlock
	LastStart Moment
	Quotum    time.Duration
}

func NewTimesheet(created time.Time) Timesheet {
	return Timesheet{
		Created:   created,
		NextId:    0,
		Blocks:    []TimeBlock{},
		LastStart: NewMoment(),
		Quotum:    0,
	}
}

func (s Timesheet) Start(t time.Time) Timesheet {
	if s.IsStarted() {
		return s
	}

	s.LastStart = NewMoment().Start(t)
	return s
}

func (s Timesheet) Stop(t time.Time, tags TagSet) Timesheet {
	if !s.IsStarted() {
		return s
	}

	newBlock := NewTimeBlock(s.NextId, s.LastStart, NewMoment().Start(t), tags)
	s.NextId += 1

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
	return s
}

func (s Timesheet) Abort() Timesheet {
	if !s.IsStarted() {
		return s
	}

	s.LastStart = NewMoment()
	return s
}

func (s Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}

func (s Timesheet) SetQuotum(q time.Duration) Timesheet {
	s.Quotum = q
	return s
}
