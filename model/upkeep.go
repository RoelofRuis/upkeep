package model

import "time"

type Upkeep struct {
	Version    string
	Categories StringStack
	Quota      map[time.Weekday]time.Duration
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

func (s Upkeep) TimesheetDuration(t Timesheet) time.Duration {
	dur := time.Duration(0)

	for _, block := range t.Blocks {
		dur += block.Duration()
	}

	if t.LastStart.IsStarted() {
		dur += time.Now().Sub(*t.LastStart.t)
	}

	return dur
}
