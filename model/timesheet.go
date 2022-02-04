package model

import (
	"time"
)

type Timesheet struct {
	// Deprecated
	Day string
	Created   time.Time
	NextId    int
	Blocks    []TimeBlock
	LastStart Moment
}

func NewTimesheet(created time.Time) *Timesheet {
	return &Timesheet{
		Created:   created,
		Blocks:    []TimeBlock{},
		LastStart: NewMoment(),
	}
}

func (s *Timesheet) Start(t time.Time) {
	if s.IsStarted() {
		return
	}

	s.LastStart = NewMoment().Start(t)
}

func (s *Timesheet) Stop(t time.Time, tags TagSet) {
	if !s.IsStarted() {
		return
	}

	newBlock := NewTimeBlock(s.NextId, s.LastStart, NewMoment().Start(t), tags)
	s.NextId += 1

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
}

func (s *Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}
