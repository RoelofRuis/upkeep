package model

import (
	"time"
)

type Timesheet struct {
	Day       string
	Blocks    []TimeBlock
	LastStart Moment
}

func NewTimesheet(day string) *Timesheet {
	return &Timesheet{
		Day:       day,
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

	newBlock := NewTimeBlock(s.LastStart, NewMoment().Start(t), tags)

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
}

func (s *Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}
