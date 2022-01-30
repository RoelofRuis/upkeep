package main

import (
	"time"
)

type Timesheet struct {
	Day    string
	Blocks []TimeBlock
}

func (s *Timesheet) TagActiveBlock(tag string) {
	if !s.HasActiveBlock() {
		return
	}

	s.Blocks[len(s.Blocks)-1].AddTag(tag)
}

func (s *Timesheet) UntagActiveBlock(tag string) {
	if !s.HasActiveBlock() {
		return
	}

	s.Blocks[len(s.Blocks)-1].RemoveTag(tag)
}

func (s *Timesheet) Start(t time.Time) {
	if s.HasActiveBlock() {
		return
	}

	s.Blocks = append(s.Blocks, NewTimeBlock(t))
}

func (s *Timesheet) Stop(t time.Time) {
	if !s.HasActiveBlock() {
		return
	}

	s.Blocks[len(s.Blocks)-1].Complete(t)
}

func (s *Timesheet) HasActiveBlock() bool {
	if len(s.Blocks) == 0 {
		return false
	}

	return !s.Blocks[len(s.Blocks)-1].HasEnded()
}
