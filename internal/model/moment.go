package model

import (
	"encoding/json"
	"errors"
	"time"
)

const LayoutDateHour = "2006-01-02 15:04 -0700 MST"
const LayoutHour = "15:04"

type Moment struct {
	t *time.Time
}

func NewMoment() Moment {
	return Moment{}
}

func (m Moment) Start(t time.Time) Moment {
	minuteRounded := t.Truncate(time.Minute)
	return Moment{t: &minuteRounded}
}

func (m Moment) IsDefined() bool {
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

func (m Moment) MarshalJSON() ([]byte, error) {
	return json.Marshal(m.Format(LayoutDateHour))
}

func (m *Moment) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case string:
		if value == "" {
			*m = NewMoment()
			return nil
		}
		t, err := time.Parse(LayoutDateHour, value)
		if err != nil {
			return err
		}
		*m = NewMoment().Start(t)
		return nil
	default:
		return errors.New("invalid time")
	}
}
