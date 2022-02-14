package model

import "time"

type BlockType = string

const (
	TypeTime     BlockType = "with-time"
	TypeDuration BlockType = "with-duration"
)

type TimeBlock struct {
	Id           int
	Category     string
	Type         BlockType
	WithTime     WithTime
	WithDuration WithDuration
}

type WithTime struct {
	Start Moment
	End   Moment
}

type WithDuration struct {
	Duration Duration
}

func NewBlockWithTime(id int, category string, start Moment, end Moment) TimeBlock {
	return TimeBlock{
		Id:       id,
		Category: category,
		Type:     TypeTime,
		WithTime: WithTime{Start: start, End: end},
	}
}

func NewBlockWithDuration(id int, category string, dur Duration) TimeBlock {
	return TimeBlock{
		Id:           id,
		Category:     category,
		Type:         TypeDuration,
		WithDuration: WithDuration{Duration: dur},
	}
}

func (b TimeBlock) HasEnded() bool {
	switch b.Type {
	case TypeTime:
		return b.WithTime.End.IsDefined()
	case TypeDuration:
		return true
	default:
		return true
	}
}

func (b TimeBlock) BaseDuration() time.Duration {
	switch b.Type {
	case TypeTime:
		return b.WithTime.End.Sub(b.WithTime.Start)
	case TypeDuration:
		return b.WithDuration.Duration.Get()
	default:
		return NewDuration().Get()
	}
}
