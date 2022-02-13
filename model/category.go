package model

type Category struct {
	Name         string
	MaxDayQuotum Duration
}

func NewCategory(name string) Category {
	return Category{
		Name:         name,
		MaxDayQuotum: NewDuration(),
	}
}
