package model

type TimeBlock struct {
	Start Moment
	End   Moment
	Tags  TagSet
}

func NewTimeBlock(start Moment, end Moment, tags TagSet) TimeBlock {
	return TimeBlock{
		Start: start,
		End:   end,
		Tags:  tags,
	}
}

func (ts *TimeBlock) AddTag(t string) {
	ts.Tags = ts.Tags.Add(t)
}

func (ts *TimeBlock) RemoveTag(t string) {
	ts.Tags = ts.Tags.Remove(t)
}
