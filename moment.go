package main

import "time"

const timeLayout = "2006-01-02 15:04"

type Moment struct {
	t *time.Time
}

func NewMoment(t time.Time) Moment {
	minuteRounded := t.Truncate(time.Minute)
	return Moment{t: &minuteRounded}
}

func NewMomentFromString(timeString string) (Moment, error) {
	if timeString == "" {
		return Moment{}, nil
	} else {
		t, err := time.Parse(timeLayout, timeString)
		if err != nil {
			return Moment{}, err
		}
		return NewMoment(t), nil
	}
}

func (m Moment) IsDefined() bool {
	return m.t != nil
}

func (m Moment) String() string {
	if m.t == nil {
		return ""
	}
	return m.t.Format(timeLayout)
}
