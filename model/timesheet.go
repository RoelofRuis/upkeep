package model

import (
	"time"
)

type Timesheet struct {
	Date      Date
	NextId    int
	Blocks    []TimeBlock
	LastStart Moment
	Quotum    Duration
}

func NewTimesheet(date Date) Timesheet {
	return Timesheet{
		Date:      date,
		NextId:    0,
		Blocks:    []TimeBlock{},
		LastStart: NewMoment(),
		Quotum:    NewDuration(),
	}
}

func (s Timesheet) Start(t time.Time) Timesheet {
	if s.IsStarted() {
		return s
	}

	s.LastStart = NewMoment().Start(t)
	return s
}

func (s Timesheet) Stop(t time.Time, category string) Timesheet {
	if !s.IsStarted() {
		return s
	}

	newBlock := NewTimeBlock(s.NextId, s.LastStart, NewMoment().Start(t), category)
	s.NextId += 1

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
	return s
}

func (s Timesheet) Remove(blockId int) Timesheet {
	for i, block := range s.Blocks {
		if block.Id == blockId {
			s.Blocks = append(s.Blocks[:i], s.Blocks[i+1:]...)
		}
	}
	return s
}

func (s Timesheet) Abort() Timesheet {
	if !s.IsStarted() {
		return s
	}

	s.LastStart = NewMoment()
	return s
}

func (s Timesheet) SetQuotum(d Duration) Timesheet {
	s.Quotum = d
	return s
}

func (s Timesheet) IsStarted() bool {
	return s.LastStart.IsStarted()
}

func (s Timesheet) GetCategoryNames() []string {
	var categories Categories
	for _, block := range s.Blocks {
		if block.Category != "" {
			categories = categories.Add(Category{Name: block.Category})
		}
	}
	return categories.Names()
}
