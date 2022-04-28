package model

import (
	"encoding/json"
	"errors"
	"time"
)

type Duration struct {
	d *time.Duration
}

func NewDuration() Duration {
	return Duration{}
}

func (d Duration) Set(dur time.Duration) Duration {
	return Duration{d: &dur}
}

func (d Duration) IsDefined() bool {
	return d.d != nil
}

func (d Duration) IsZero() bool {
	return !d.IsDefined() || *d.d == 0
}

func (d Duration) AddDuration(dur time.Duration) Duration {
	if d.IsDefined() {
		return d.Add(NewDuration().Set(dur))
	}
	return NewDuration().Set(dur)
}

func (d Duration) Add(other Duration) Duration {
	if d.IsDefined() {
		if other.IsDefined() {
			return NewDuration().Set(d.Get() + other.Get())
		}
		return d
	}
	return other
}

func (d Duration) Sub(other Duration) Duration {
	if d.IsDefined() {
		if other.IsDefined() {
			return NewDuration().Set(d.Get() - other.Get())
		}
		return d
	}
	return other
}

func (d Duration) Get() time.Duration {
	if !d.IsDefined() {
		return time.Duration(0)
	}
	return *d.d
}

func (d Duration) MarshalJSON() ([]byte, error) {
	if d.d == nil {
		return json.Marshal("")
	}

	return json.Marshal(d.d.String())
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = NewDuration().Set(time.Duration(value))
		return nil
	case string:
		if value == "" {
			*d = NewDuration()
			return nil
		}
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = NewDuration().Set(dur)
	default:
		return errors.New("Invalid duration")
	}
	return nil
}
