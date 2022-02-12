package model

import "time"

type TimeBlock struct {
	Id       int
	Category string
	Start    Moment
	End      Moment
}

type TimedBlock struct {
	Start Moment
	End   Moment
}

type DurationBlock struct {
	Duration time.Duration
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
