package model

import "time"

type Upkeep struct {
	Version string
	Tags    TagStack
	Quota   map[time.Weekday]time.Duration
}

func (s Upkeep) ShiftTags() Upkeep {
	s.Tags = s.Tags.Push(NewTagSet())
	return s
}

func (s Upkeep) UnshiftTags() Upkeep {
	stack, _, _ := s.Tags.Pop()
	s.Tags = stack
	return s
}

func (s *Upkeep) GetTags() TagSet {
	set, has := s.Tags.Peek()
	if !has {
		return NewTagSet()
	}
	return set
}

func (s Upkeep) AddTag(tag string) Upkeep {
	stack, set, has := s.Tags.Pop()
	if !has {
		set = NewTagSet()
	}
	s.Tags = stack.Push(set.Add(tag))

	return s
}

func (s Upkeep) RemoveTag(tag string) Upkeep {
	stack, set, has := s.Tags.Pop()
	if !has {
		return s
	}
	s.Tags = stack.Push(set.Remove(tag))

	return s
}

func (s Upkeep) SetQuotumForDay(day time.Weekday, quotum time.Duration) Upkeep {
	s.Quota[day] = quotum
	return s
}

func (s Upkeep) GetQuotumForDay(day time.Weekday) time.Duration {
	quotum, has := s.Quota[day]
	if !has {
		return 0
	}

	return quotum
}