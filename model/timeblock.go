package model

import "time"

type TimeBlock struct {
	Id       int
	Start    Moment
	End      Moment
	Category string
}

func NewTimeBlock(id int, start Moment, end Moment, category string) TimeBlock {
	return TimeBlock{
		Id:       id,
		Start:    start,
		End:      end,
		Category: category,
	}
}

func (b TimeBlock) BaseDuration() time.Duration {
	return b.End.Sub(b.Start)
}
