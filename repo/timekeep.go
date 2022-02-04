package repo

import (
	"timesheet/infra"
	"timesheet/model"
)

type UpkeepRepository struct {
	FileIO infra.FileIO
}

type timekeepJson struct {
	Version string `json:"version"`
	Tags    string `json:"tags"`
}

func (r *UpkeepRepository) Get() (*model.Upkeep, error) {
	input := timekeepJson{
		Version: "0.1",
	}

	if err := r.FileIO.Read("timekeep.json", &input); err != nil {
		return nil, err
	}

	timekeep := &model.Upkeep{
		Version: input.Version,
		Tags:    model.NewTagStackFromString(input.Tags),
	}

	return timekeep, nil
}

func (r *UpkeepRepository) Insert(m *model.Upkeep) error {
	output := timekeepJson{
		Version: m.Version,
		Tags:    m.Tags.String(),
	}

	if err := r.FileIO.Write("timekeep.json", output); err != nil {
		return err
	}

	return nil
}
