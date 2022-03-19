package repo

import (
	"github.com/roelofruis/upkeep/internal/infra"
)

type Repository struct {
	Upkeep     UpkeepRepository
	Timesheets TimesheetRepository
}

func New(io infra.IO) Repository {
	return Repository{
		Upkeep:     UpkeepRepository{io},
		Timesheets: TimesheetRepository{io},
	}
}
