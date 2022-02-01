package model

import (
	"time"
)

type Timesheet struct {
	Day       string
	LastStart Moment
	Blocks    []TimeBlock
}

func NewTimesheet(day string) *Timesheet {
	return &Timesheet{
		Day:       day,
		LastStart: NewMoment(),
		Blocks:    []TimeBlock{},
	}
}

func (s *Timesheet) TagActiveBlock(tag string) {
	if !s.IsStarted() {
		return
	}

	s.Blocks[len(s.Blocks)-1].AddTag(tag)
}

func (s *Timesheet) UntagActiveBlock(tag string) {
	if !s.IsStarted() {
		return
	}

	s.Blocks[len(s.Blocks)-1].RemoveTag(tag)
}

func (s *Timesheet) Start(t time.Time) {
	if s.IsStarted() {
		return
	}

	s.LastStart = NewMoment().Start(t)
}

func (s *Timesheet) Stop(t time.Time) {
	if !s.IsStarted() {
		return
	}

	s.Blocks = append(s.Blocks, NewTimeBlock(s.LastStart, NewMoment().Start(t)))
	s.LastStart = NewMoment()
}

func (s *Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}
