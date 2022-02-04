package repo

import (
	"timesheet/infra"
	"timesheet/model"
)

type TimekeepRepository struct {
	FileIO infra.FileIO
}

type timekeepJson struct {
	Version string `json:"version"`
	Tags    string `json:"tags"`
}

func (r *TimekeepRepository) Get() (*model.Timekeep, error) {
	input := timekeepJson{
		Version: "0.1",
	}

	if err := r.FileIO.Read("timekeep.json", &input); err != nil {
		return nil, err
	}

	timekeep := &model.Timekeep{
		Version: input.Version,
		Tags:    model.NewTagStackFromString(input.Tags),
	}

	return timekeep, nil
}

func (r *TimekeepRepository) Insert(m *model.Timekeep) error {
	output := timekeepJson{
		Version: m.Version,
		Tags:    m.Tags.String(),
	}

	if err := r.FileIO.Write("timekeep.json", output); err != nil {
		return err
	}

	return nil
}
