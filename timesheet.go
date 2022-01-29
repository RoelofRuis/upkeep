package main

import (
	"errors"
	"time"
)

type Timesheet struct {
	Blocks []TimeBlock
}

func (s *Timesheet) Start(t time.Time) error {
	if s.HasActiveBlock() {
		return errors.New("cannot start two blocks")
	}

	s.Blocks = append(s.Blocks, NewTimeBlock(t))
	return nil
}

func (s *Timesheet) Stop(t time.Time) error {
	if !s.HasActiveBlock() {
		return errors.New("no block active")
	}

	s.Blocks[len(s.Blocks)-1].Complete(t)
	return nil
}

func (s *Timesheet) HasActiveBlock() bool {
	if len(s.Blocks) == 0 {
		return false
	}

	return !s.Blocks[len(s.Blocks)-1].HasEnded()
}
