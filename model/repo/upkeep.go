package repo

import (
	"time"
	"upkeep/infra"
	"upkeep/model"
)

type UpkeepRepository struct {
	FileIO infra.FileIO
}

type upkeepJson struct {
	Version            string         `json:"version"`
	ActiveCategory     string         `json:"active_category"`
	Quota              map[int]string `json:"quota"`
	ExcludedCategories string         `json:"excluded_categories"`
}

const VERSION = "0.2"

func (r *UpkeepRepository) Get() (model.Upkeep, error) {
	input := upkeepJson{
		Version: "0.1",
	}

	if err := r.FileIO.Read("upkeep.json", &input); err != nil {
		return model.Upkeep{}, err
	}

	quotumMap := make(map[time.Weekday]time.Duration)
	for weekday, dur := range input.Quota {
		duration, err := time.ParseDuration(dur)
		if err != nil {
			return model.Upkeep{}, err
		}
		quotumMap[time.Weekday(weekday)] = duration
	}

	upkeep := model.Upkeep{
		Version:            input.Version,
		Categories:         infra.NewStackFromString(input.ActiveCategory),
		Quota:              quotumMap,
		ExcludedCategories: infra.NewSetFromString(input.ExcludedCategories),
	}

	upkeep.Version = VERSION

	return upkeep, nil
}

func (r *UpkeepRepository) Insert(m model.Upkeep) error {
	quotumMap := make(map[int]string)
	for weekday, dur := range m.Quota {
		quotumMap[int(weekday)] = dur.String()
	}

	output := upkeepJson{
		Version:            m.Version,
		ActiveCategory:     m.Categories.String(),
		Quota:              quotumMap,
		ExcludedCategories: m.ExcludedCategories.String(),
	}

	if err := r.FileIO.Write("upkeep.json", output); err != nil {
		return err
	}

	return nil
}
