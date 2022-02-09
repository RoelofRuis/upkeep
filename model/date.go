package model

import (
	"time"
)

type Date time.Time

func NewDateFromString(s string) (Date, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return NewDate(t), nil
}

func NewDate(t time.Time) Date {
	return Date(time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC))
}

func Today() Date {
	return NewDate(time.Now())
}

func (d Date) DayBefore() Date {
	return Date(time.Time(d).AddDate(0, 0, -1))
}

func (d Date) Format(layout string) string {
	return time.Time(d).Format(layout)
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d Date) IsToday() bool {
	return time.Time(d).Equal(time.Time(Today()))
}