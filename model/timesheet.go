package model

import (
	"time"
)

type Timesheet struct {
	Day    string
	Break  bool // TODO: does this belong here?
	Blocks []TimeBlock
	// Deprecated should be moved to Timekeep
	LastStart Moment
	// Deprecated should be moved to Timekeep
	Tags TagSet
}

func NewTimesheet(day string) *Timesheet { // TODO: remove day from constructor?
	return &Timesheet{
		Day:       day,
		Break:     false,
		LastStart: NewMoment(),
		Tags:      NewTagSet(),
		Blocks:    []TimeBlock{},
	}
}

func (s *Timesheet) SetBreak(active bool) {
	s.Break = active
}

func (s *Timesheet) AttachTag(tag string) {
	s.Tags = s.Tags.Add(tag)
}

func (s *Timesheet) DetachTag(tag string) {
	s.Tags = s.Tags.Remove(tag)
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

	var tags TagSet
	if !s.Break {
		tags = s.Tags
	} else {
		tags = NewTagSetFromString("break")
	}

	newBlock := NewTimeBlock(s.LastStart, NewMoment().Start(t), tags)

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
}

func (s *Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}
