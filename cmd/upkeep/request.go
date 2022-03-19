package main

import (
	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

type Request struct {
	Clock      infra.Clock
	Params     infra.Params
	Upkeep     *model.Upkeep
	Timesheets []*model.Timesheet
}
