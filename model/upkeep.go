package model

import (
	"time"
	"upkeep/infra"
)

type Upkeep struct {
	Version            string
	Categories         infra.Stack
	Quota              map[time.Weekday]time.Duration
	Discounts          []Discount
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

func (s Upkeep) RemoveDiscount(category string) Upkeep {
	for i, d := range s.Discounts {
		if d.Category == category {
			s.Discounts[i] = s.Discounts[len(s.Discounts)-1]
			s.Discounts = s.Discounts[:len(s.Discounts)-1]
		}
	}
	return s
}

func (s Upkeep) SetDiscount(d Discount) Upkeep {
	s.Discounts = append(s.Discounts, d)
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

func (s Upkeep) DiscountApplies(cat string) bool {
	for _, d := range s.Discounts {
		if d.Category == cat {
			return true
		}
	}
	return false
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

	var discountMeasures []*Discounter
	for _, d := range s.Discounts {
		discountMeasures = append(discountMeasures, d.Measure())
	}

	for _, block := range t.Blocks {
		blockDur := block.Duration()
		for _, m := range discountMeasures {
			blockDur = m.GetTimeRemaining(block.Category, blockDur)
		}
		dur += blockDur
	}

	if t.LastStart.IsStarted() && t.Date.IsToday() {
		blockDur := time.Now().Sub(*t.LastStart.t)
		for _, m := range discountMeasures {
			blockDur = m.GetTimeRemaining(s.Categories.Peek(), blockDur)
		}
		dur += blockDur
	}

	return dur
}
