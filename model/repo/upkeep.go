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
	Version    string         `json:"version"`
	Categories string         `json:"categories"`
	Quota      map[int]string `json:"quota"`
}

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
		Version:    input.Version,
		Categories: model.NewStringStackFromString(input.Categories),
		Quota:      quotumMap,
	}

	return upkeep, nil
}

func (r *UpkeepRepository) Insert(m model.Upkeep) error {
	quotumMap := make(map[int]string)
	for weekday, dur := range m.Quota {
		quotumMap[int(weekday)] = dur.String()
	}

	output := upkeepJson{
		Version:    m.Version,
		Categories: m.Categories.String(),
		Quota:      quotumMap,
	}

	if err := r.FileIO.Write("upkeep.json", output); err != nil {
		return err
	}

	return nil
}
