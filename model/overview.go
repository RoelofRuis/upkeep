package model

import "time"

type DiscountedTimeBlocks []DiscountedTimeBlock

func (d DiscountedTimeBlocks) TotalDuration() time.Duration {
	dur := time.Duration(0)
	for _, b := range d {
		dur += b.DiscountedDuration
	}
	return dur
}

type DiscountedTimeBlock struct {
	Block              TimeBlock
	IsDiscounted       bool
	DiscountedDuration time.Duration
}
