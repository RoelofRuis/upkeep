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
	Finalised bool
}

func NewTimesheet(date Date) Timesheet {
	return Timesheet{
		Date:      date,
		NextId:    0,
		Blocks:    []TimeBlock{},
		LastStart: NewMoment(),
		Quotum:    NewDuration(),
		Finalised: false,
	}
}

func (s Timesheet) Start(t time.Time) Timesheet {
	if s.Finalised || s.IsStarted() {
		return s
	}

	s.LastStart = NewMoment().Start(t)
	return s
}

func (s Timesheet) Stop(t time.Time, category string) Timesheet {
	if s.Finalised || !s.IsStarted() {
		return s
	}

	newBlock := NewBlockWithTime(s.NextId, category, false, s.LastStart, NewMoment().Start(t))
	s.NextId += 1

	s.Blocks = append(s.Blocks, newBlock)
	s.LastStart = NewMoment()
	return s
}

func (s Timesheet) Write(category string, dur Duration) Timesheet {
	if s.Finalised {
		return s
	}

	newBlock := NewBlockWithDuration(s.NextId, category, false, dur)
	s.NextId += 1

	s.Blocks = append(s.Blocks, newBlock)
	return s
}

func (s Timesheet) UpdateBlockCategory(blockId int, category string) Timesheet {
	for i, block := range s.Blocks {
		if block.Id == blockId {
			s.Blocks[i].Category = category
		}
	}
	return s
}

func (s Timesheet) RemoveBlock(blockId int) Timesheet {
	if s.Finalised {
		return s
	}

	for i, block := range s.Blocks {
		if block.Id == blockId {
			s.Blocks[i].Deleted = true
			break
		}
	}
	return s
}

func (s Timesheet) RestoreBlock(blockId int) Timesheet {
	if s.Finalised {
		return s
	}

	for i, block := range s.Blocks {
		if block.Id == blockId {
			s.Blocks[i].Deleted = false
			break
		}
	}
	return s
}

func (s Timesheet) Abort() Timesheet {
	if s.Finalised || !s.IsStarted() {
		return s
	}

	s.LastStart = NewMoment()
	return s
}

func (s Timesheet) Finalise() Timesheet {
	aborted := s.Abort()
	aborted.Finalised = true
	return aborted
}

func (s Timesheet) Unfinalise() Timesheet {
	s.Finalised = false
	return s
}

func (s Timesheet) SetQuotum(d Duration) Timesheet {
	s.Quotum = d
	return s
}

func (s Timesheet) IsStarted() bool {
	return s.LastStart.IsDefined()
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
