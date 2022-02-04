package repo

import (
	"timesheet/infra"
	"timesheet/model"
)

type UpkeepRepository struct {
	FileIO infra.FileIO
}

type upkeepJson struct {
	Version string `json:"version"`
	Tags    string `json:"tags"`
}

func (r *UpkeepRepository) Get() (*model.Upkeep, error) {
	input := upkeepJson{
		Version: "0.1",
	}

	if err := r.FileIO.Read("upkeep.json", &input); err != nil {
		return nil, err
	}

	upkeep := &model.Upkeep{
		Version: input.Version,
		Tags:    model.NewTagStackFromString(input.Tags),
	}

	return upkeep, nil
}

func (r *UpkeepRepository) Insert(m *model.Upkeep) error {
	output := upkeepJson{
		Version: m.Version,
		Tags:    m.Tags.String(),
	}

	if err := r.FileIO.Write("upkeep.json", output); err != nil {
		return err
	}

	return nil
}
