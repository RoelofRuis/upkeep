package model

import (
	"time"
	"upkeep/infra"
)

type Upkeep struct {
	Version            string
	SelectedCategories infra.Stack
	Quota              map[time.Weekday]time.Duration
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
	cat.MaxDayQuotum = dur
	s.Categories = s.Categories.Add(cat)
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

func (s Upkeep) DiscountTimeBlocks(t Timesheet) DiscountedTimeBlocks {
	categoryQuota := make(map[string]time.Duration)
	for _, c := range s.Categories {
		if c.MaxDayQuotum != nil {
			categoryQuota[c.Name] = *c.MaxDayQuotum
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
		if t.Date.IsToday() {
			discountedDur := time.Now().Sub(*t.LastStart.t)
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
					End:      NewMoment().Start(time.Now()),
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
