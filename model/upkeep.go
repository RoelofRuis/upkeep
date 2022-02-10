package model

import (
	"time"
	"upkeep/infra"
)

type Upkeep struct {
	Version            string
	Categories         infra.Stack
	Quota              map[time.Weekday]time.Duration
	ExcludedCategories infra.Set
}

func (s Upkeep) ShiftCategory() Upkeep {
	s.Categories = s.Categories.Push("")
	return s
}

func (s Upkeep) UnshiftCategory() Upkeep {
	stack, _, _ := s.Categories.Pop()
	s.Categories = stack
	return s
}

func (s *Upkeep) GetCategory() string {
	return s.Categories.Peek()
}

func (s Upkeep) SetCategory(name string) Upkeep {
	stack, _, _ := s.Categories.Pop()
	s.Categories = stack.Push(name)

	return s
}

func (s Upkeep) AddExcludedCategory(name string) Upkeep {
	s.ExcludedCategories = s.ExcludedCategories.Add(name)

	return s
}

func (s Upkeep) RemoveExcludedCategory(name string) Upkeep {
	s.ExcludedCategories = s.ExcludedCategories.Remove(name)

	return s
}

func (s Upkeep) SetQuotumForDay(day time.Weekday, quotum time.Duration) Upkeep {
	s.Quota[day] = quotum
	return s
}

func (s Upkeep) RemoveQuotumForDay(day time.Weekday) Upkeep {
	delete(s.Quota, day)
	return s
}

func (s Upkeep) GetQuotumForDay(day time.Weekday) time.Duration {
	quotum, has := s.Quota[day]
	if !has {
		return 0
	}

	return quotum
}

func (s Upkeep) TimesheetQuotum(t Timesheet) time.Duration {
	quotum := t.Quotum

	if quotum == 0 {
		weekdayQuotum, has := s.Quota[t.Date.Weekday()]
		if !has {
			return 0
		}
		return weekdayQuotum
	}

	return t.Quotum
}

func (s Upkeep) TimesheetDuration(t Timesheet) time.Duration {
	dur := time.Duration(0)

	for _, block := range t.Blocks {
		if !s.ExcludedCategories.Contains(block.Category) {
			dur += block.Duration()
		}
	}

	if t.LastStart.IsStarted() && t.Date.IsToday() && !s.ExcludedCategories.Contains(s.Categories.Peek()) {
		dur += time.Now().Sub(*t.LastStart.t)
	}

	return dur
}
