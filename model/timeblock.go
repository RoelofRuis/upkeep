package model

import "time"

type TimeBlock struct {
	Id    int
	Start Moment
	End   Moment
	Tags  TagSet
}

func NewTimeBlock(id int, start Moment, end Moment, tags TagSet) TimeBlock {
	return TimeBlock{
		Id:    id,
		Start: start,
		End:   end,
		Tags:  tags,
	}
}

func (b TimeBlock) Duration() time.Duration {
	return b.End.t.Sub(*b.Start.t)
}
