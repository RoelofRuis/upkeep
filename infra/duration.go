package infra

import (
	"encoding/json"
	"errors"
	"time"
)

type JSONDuration time.Duration

func NewDuration() JSONDuration {
	return JSONDuration(0)
}

func (d JSONDuration) Unpack() time.Duration {
	return time.Duration(d)
}

func (d JSONDuration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *JSONDuration) UnmarshalJSON(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = JSONDuration(value)
		return nil
	case string:
		dur, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = JSONDuration(dur)
	default:
		return errors.New("Invalid duration")
	}
	return nil
}
