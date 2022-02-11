package model

import "time"

type Category struct {
	Name         string
	MaxDayQuotum *time.Duration
}

func NewCategory(name string) Category {
	return Category{
		Name:         name,
		MaxDayQuotum: nil,
	}
}
