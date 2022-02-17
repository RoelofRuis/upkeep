package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

func (d Date) ShiftDay(days int) Date {
	return Date(time.Time(d).AddDate(0, 0, days))
}

func (d Date) PreviousMonday() Date {
	daysBack := int((time.Time(d).Weekday() + 6) % 7)
	return Date(time.Time(d).AddDate(0, 0, -daysBack))
}

func (d Date) FirstOfMonth() Date {
	t := time.Time(d)
	return Date(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC))
}

func (d Date) DaysInMonth() int {
	t := time.Time(d)
	return time.Date(t.Year(), t.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func (d Date) IterateNext(num int) []Date {
	dates := make([]Date, num)
	for i := 0; i < num; i++ {
		dates[i] = d.ShiftDay(i)
	}
	return dates
}

func (d Date) Weekday() time.Weekday {
	return time.Time(d).Weekday()
}

func (d Date) Format(layout string) string {
	return time.Time(d).Format(layout)
}

func (d Date) String() string {
	return d.Format("2006-01-02")
}

func (d Date) OnSameDateAs(t time.Time) bool {
	return time.Time(d).Equal(time.Time(NewDate(t)))
}

func (d Date) Year() int {
	return time.Time(d).Year()
}

func (d Date) Month() time.Month {
	return time.Time(d).Month()
}

func (d Date) Day() int {
	return time.Time(d).Day()
}

func (d Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Format("2006-01-02"))
}

func (d *Date) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		date, err := time.Parse("2006-01-02", value)
		if err != nil {
			return err
		}
		*d = Date(date)
	default:
		return errors.New("invalid date")
	}
	return nil
}

var quickDateRegex = regexp.MustCompile("^(-?[0-9]+)?([a-z]+)$")

func IterDates(dateDef string) ([]Date, error) {
	date := NewDate(time.Now())
	shifts := 0
	numDays := 1

	matches := quickDateRegex.FindStringSubmatch(dateDef)
	if len(matches) == 3 {
		if matches[1] != "" {
			i, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return nil, err
			}
			shifts = int(i)
		}
		dateDef = matches[2]
	}

	switch dateDef {
	case "day":
	case "d":
		date = date.ShiftDay(shifts)
		break

	case "week":
	case "w":
		date = NewDate(time.Now()).PreviousMonday().ShiftDay(shifts * 7)
		numDays = 5
		break

	case "weekfull":
	case "wf":
		date = NewDate(time.Now()).PreviousMonday().ShiftDay(shifts * 7)
		numDays = 7
		break

	default:
		parsedDate, err := NewDateFromString(dateDef)
		if err != nil {
			return nil, fmt.Errorf("invalid date value '%s'", dateDef)
		}
		date = parsedDate
		break
	}

	return date.IterateNext(numDays), nil
}
