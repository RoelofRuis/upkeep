package model

import (
	"time"
	"upkeep/infra"
)

type Upkeep struct {
	Version            string
	SelectedCategories infra.Stack
	Quota              map[time.Weekday]Duration
	Categories         Categories
}

func (s Upkeep) ShiftSelectedCategory() Upkeep {
	s.SelectedCategories = s.SelectedCategories.Push("")
	return s
}

func (s Upkeep) UnshiftSelectedCategory() Upkeep {
	stack, _, _ := s.SelectedCategories.Pop()
	s.SelectedCategories = stack
	return s
}

func (s *Upkeep) GetSelectedCategory() Category {
	selected := s.SelectedCategories.Peek()
	return s.Categories.Get(selected)
}

func (s Upkeep) SetSelectedCategory(name string) Upkeep {
	stack, _, _ := s.SelectedCategories.Pop()
	s.SelectedCategories = stack.Push(name)

	return s
}

func (s Upkeep) SetCategoryMaxDayQuotum(category string, dur *time.Duration) Upkeep {
	cat := s.Categories.Get(category)

	newDur := NewDuration()
	if dur != nil {
		newDur = newDur.Set(*dur)
	}
	cat.MaxDayQuotum = newDur

	s.Categories = s.Categories.Add(cat)
	return s
}

func (s Upkeep) SetQuotumForDay(day time.Weekday, quotum time.Duration) Upkeep {
	s.Quota[day] = NewDuration().Set(quotum)
	return s
}

func (s Upkeep) RemoveQuotumForDay(day time.Weekday) Upkeep {
	delete(s.Quota, day)
	return s
}

func (s Upkeep) GetQuotumForDay(day time.Weekday) Duration {
	quotum, has := s.Quota[day]
	if !has {
		return NewDuration()
	}

	return quotum
}

func (s Upkeep) TimesheetQuotum(t Timesheet) Duration {
	quotum := t.Quotum

	if !quotum.IsDefined() {
		weekdayQuotum, has := s.Quota[t.Date.Weekday()]
		if !has {
			return NewDuration()
		}
		return weekdayQuotum
	}

	return t.Quotum
}

func (s Upkeep) DiscountTimeBlocks(t Timesheet, at time.Time) DiscountedTimeBlocks {
	categoryQuota := make(map[string]time.Duration)
	for _, c := range s.Categories {
		if c.MaxDayQuotum.IsDefined() {
			categoryQuota[c.Name] = *c.MaxDayQuotum.d
		}
	}

	var discountedBlocks []DiscountedTimeBlock

	for _, block := range t.Blocks {
		discountedDur := block.BaseDuration()
		isDiscounted := false
		remaining, has := categoryQuota[block.Category]
		if has {
			if remaining > discountedDur {
				categoryQuota[block.Category] -= discountedDur
			} else {
				discountedDur = remaining
				categoryQuota[block.Category] = 0
				isDiscounted = true
			}
		}
		discountedBlocks = append(discountedBlocks, DiscountedTimeBlock{
			Block:              block,
			IsDiscounted:       isDiscounted,
			DiscountedDuration: discountedDur,
		})
	}

	if t.LastStart.IsStarted() {
		cat := s.SelectedCategories.Peek()
		if t.Date.OnSameDateAs(at) {
			discountedDur := at.Sub(*t.LastStart.t)
			isDiscounted := false
			remaining, has := categoryQuota[cat]
			if has {
				if remaining > discountedDur {
					categoryQuota[cat] -= discountedDur
				} else {
					discountedDur = remaining
					categoryQuota[cat] = 0
					isDiscounted = true
				}
			}
			discountedBlocks = append(discountedBlocks, DiscountedTimeBlock{
				Block: TimeBlock{
					Id:       -1,
					Category: cat,
					Start:    t.LastStart,
					End:      NewMoment().Start(at),
				},
				IsDiscounted:       isDiscounted,
				DiscountedDuration: discountedDur,
			})
		} else {
			discountedBlocks = append(discountedBlocks, DiscountedTimeBlock{
				Block: TimeBlock{
					Id:       -1,
					Category: cat,
					Start:    t.LastStart,
					End:      NewMoment(),
				},
				IsDiscounted:       false,
				DiscountedDuration: 0,
			})
		}
	}

	return discountedBlocks
}
