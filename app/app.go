package app

import (
	"time"
	"upkeep/infra"
	"upkeep/model"
)

type App struct {
	Params     infra.Params
	Upkeep     *model.Upkeep
	Timesheets []*model.Timesheet
}

func (r Repository) Handle(f func(editor *App) (string, error)) infra.Handler {
	return func(params infra.Params) (string, error) {
		upkeep, err := r.Upkeep.Get()
		if err != nil {
			return "", err
		}

		date, numDays, err := MakeDateRange(model.NewDate(time.Now()), params)
		if err != nil {
			return "", err
		}

		dates := date.IterateNext(numDays)
		refSheets := make([]*model.Timesheet, len(dates))
		timesheets := make([]*model.Timesheet, len(dates))
		for i, day := range dates {
			sheet, err := r.Timesheets.GetForDate(day)
			if err != nil {
				return "", err
			}
			timesheets[i] = &sheet
			refSheets[i] = &sheet
		}

		app := &App{
			Params:     params,
			Upkeep:     &upkeep,
			Timesheets: timesheets,
		}

		s, err := f(app)
		if err != nil {
			return s, err
		}

		if app.Upkeep != &upkeep {
			if err := r.Upkeep.Insert(*app.Upkeep); err != nil {
				return "", err
			}
		}
		for i := 0; i < len(refSheets); i++ {
			if app.Timesheets[i] != refSheets[i] {
				if err := r.Timesheets.Insert(*app.Timesheets[i]); err != nil {
					return "", err
				}
			}
		}

		return s, nil
	}
}
