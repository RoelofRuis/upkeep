package model

import (
	"time"

	"github.com/roelofruis/upkeep/internal/infra"
)

type Upkeep struct {
	Version            string
	SelectedCategories infra.Stack
	Quota              map[time.Weekday]Duration
	CategorySettings   CategorySettings
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

func (s Upkeep) SwapCategories() Upkeep {
	lastStack, lastCat, has := s.SelectedCategories.Pop()
	if !has {
		return s
	}
	stack, secondLastCat, has := lastStack.Pop()
	if !has {
		return s
	}
	s.SelectedCategories = stack.Push(lastCat).Push(secondLastCat)
	return s
}

func (s *Upkeep) GetSelectedCategory() CategorySetting {
	selected := s.SelectedCategories.Peek()
	return s.CategorySettings.Get(selected)
}

func (s Upkeep) SetSelectedCategory(name string) Upkeep {
	stack, _, _ := s.SelectedCategories.Pop()
	s.SelectedCategories = stack.Push(name)

	return s
}

func (s Upkeep) SetCategoryMaxDayQuotum(category string, dur *time.Duration) Upkeep {
	cat := s.CategorySettings.Get(category)

	newDur := NewDuration()
	if dur != nil {
		newDur = newDur.Set(*dur)
	}
	cat.MaxDayQuotum = newDur

	s.CategorySettings = s.CategorySettings.Add(cat)
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

func (s Upkeep) GetWeekdayQuotum(day time.Weekday) Duration {
	quotum, has := s.Quota[day]
	if !has {
		return NewDuration()
	}

	return quotum
}

func (s Upkeep) GetTimesheetQuotum(t Timesheet) Duration {
	quotum := t.Quotum

	if quotum.IsZero() {
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
	for _, c := range s.CategorySettings {
		if c.MaxDayQuotum.IsDefined() {
			categoryQuota[c.Name] = *c.MaxDayQuotum.d
		}
	}

	var discountedBlocks []DiscountedTimeBlock

	for _, block := range t.Blocks {
		if block.Deleted {
			continue
		}
		discountedDur := block.BaseDuration()
		isDiscounted := false
		remaining, has := categoryQuota[block.Category.String()]
		if has {
			if remaining > discountedDur {
				categoryQuota[block.Category.String()] -= discountedDur
			} else {
				discountedDur = remaining
				categoryQuota[block.Category.String()] = 0
				isDiscounted = true
			}
		}
		discountedBlocks = append(discountedBlocks, DiscountedTimeBlock{
			Block:              block,
			IsDiscounted:       isDiscounted,
			DiscountedDuration: discountedDur,
		})
	}

	if t.LastStart.IsDefined() {
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
				Block:              NewBlockWithTime(-1, cat, false, t.LastStart, NewMoment().Start(at)),
				IsDiscounted:       isDiscounted,
				DiscountedDuration: discountedDur,
			})
		} else {
			discountedBlocks = append(discountedBlocks, DiscountedTimeBlock{
				Block:              NewBlockWithTime(-1, cat, false, t.LastStart, NewMoment()),
				IsDiscounted:       false,
				DiscountedDuration: 0,
			})
		}
	}

	return discountedBlocks
}
