package repo

import (
	"github.com/roelofruis/upkeep/internal/infra"
)

type Repository struct {
	Upkeep     UpkeepRepository
	Timesheets TimesheetRepository
}

func New(fileIO infra.FileIO) Repository {
	return Repository{
		Upkeep:     UpkeepRepository{fileIO},
		Timesheets: TimesheetRepository{fileIO},
	}
}
