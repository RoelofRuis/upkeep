package model

type Upkeep struct {
	Version string
	Tags    TagStack
}

func (s *Upkeep) ShiftTags() {
	s.Tags = s.Tags.Push(NewTagSet())
}

func (s *Upkeep) UnshiftTags() {
	stack, _, _ := s.Tags.Pop()
	s.Tags = stack
}

func (s *Upkeep) GetTags() TagSet {
	set, has := s.Tags.Peek()
	if !has {
		return NewTagSet()
	}
	return set
}

func (s *Upkeep) AddTag(tag string) {
	stack, set, has := s.Tags.Pop()
	if !has {
		set = NewTagSet()
	}
	s.Tags = stack.Push(set.Add(tag))
}

func (s *Upkeep) RemoveTag(tag string) {
	stack, set, has := s.Tags.Pop()
	if !has {
		return
	}
	s.Tags = stack.Push(set.Remove(tag))
}