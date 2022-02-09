package model

import "time"

const LayoutDateHour = "2006-01-02 15:04 -0700 MST"
const LayoutHour = "15:04"

type Moment struct {
	t *time.Time
}

func NewMoment() Moment {
	return Moment{}
}

func NewMomentFromString(timeString string) (Moment, error) {
	if timeString == "" {
		return NewMoment(), nil
	} else {
		t, err := time.Parse(LayoutDateHour, timeString)
		if err != nil {
			return Moment{}, err
		}
		return NewMoment().Start(t), nil
	}
}

func (m Moment) Start(t time.Time) Moment {
	minuteRounded := t.Truncate(time.Minute)
	return Moment{t: &minuteRounded}
}

func (m Moment) IsStarted() bool {
	return m.t != nil
}

func (m Moment) Format(layout string) string {
	if m.t == nil {
		return ""
	}
	return m.t.Format(layout)
}

func (m Moment) Sub(that Moment) time.Duration {
	if m.t == nil || that.t == nil {
		return 0
	}

	return m.t.Sub(*that.t)
}