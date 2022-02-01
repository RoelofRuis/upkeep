package model

import (
	"time"
)

type Timesheet struct {
	Day       string
	LastStart Moment
	Tags      TagSet
	Blocks    []TimeBlock
}

func NewTimesheet(day string) *Timesheet {
	return &Timesheet{
		Day:       day,
		LastStart: NewMoment(),
		Tags:      NewTagSet(),
		Blocks:    []TimeBlock{},
	}
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

	newBlock := NewTimeBlock(s.LastStart, NewMoment().Start(t), s.Tags)

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
}

func (s *Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}
