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
	Version string         `json:"version"`
	Tags    string         `json:"tags"`
	Quota   map[int]string `json:"quota"`
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
		Version: input.Version,
		Tags:    model.NewTagStackFromString(input.Tags),
		Quota:   quotumMap,
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
		Version: m.Version,
		Tags:    m.Tags.String(),
		Quota:   quotumMap,
	}

	if err := r.FileIO.Write("upkeep.json", output); err != nil {
		return err
	}

	return nil
}
