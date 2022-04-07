package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/roelofruis/upkeep/internal/model"
	"github.com/roelofruis/upkeep/internal/model/repo"
)

type Repository repo.Repository

func HandleStart(app *App) (string, error) {
	category := app.Params.Get(0)

	_, err := HandleStop(app)
	if err != nil {
		return "", err
	}

	now := time.Now()
	sheet := app.Timesheets[0].Start(now)

	quotum := app.Upkeep.GetWeekdayQuotum(now.Weekday())
	sheet = sheet.SetQuotum(quotum)

	app.Timesheets[0] = &sheet

	if category != "" {
		_, err := HandleSet(app)
		if err != nil {
			return "", err
		}
	}

	return ViewSheets(app)
}

func HandleStop(app *App) (string, error) {
	sheet := app.Timesheets[0].Stop(time.Now(), app.Upkeep.GetSelectedCategory().Name)
	app.Timesheets[0] = &sheet

	return ViewSheets(app)
}

func HandleAbort(app *App) (string, error) {
	sheet := app.Timesheets[0].Abort()
	app.Timesheets[0] = &sheet

	return ViewSheets(app)
}

func HandleSwitch(app *App) (string, error) {
	_, err := HandleStop(app)
	if err != nil {
		return "", err
	}

	upkeep := app.Upkeep.ShiftSelectedCategory()
	app.Upkeep = &upkeep

	return HandleStart(app)
}

func HandleSwap(app *App) (string, error) {
	upkeep := app.Upkeep.SwapCategories()
	app.Upkeep = &upkeep

	return ViewSheets(app)
}

func HandleContinue(app *App) (string, error) {
	_, err := HandleStop(app)
	if err != nil {
		return "", err
	}

	upkeep := app.Upkeep.UnshiftSelectedCategory()
	app.Upkeep = &upkeep
	return HandleStart(app)
}

var validCategory = regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)?$`)

func HandleSet(app *App) (string, error) {
	if app.Params.Len() == 0 {
		return "", errors.New("no category specified")
	}

	upkeep := *app.Upkeep
	category := app.Params.Get(0)

	if !validCategory.MatchString(category) {
		return ViewSheets(app)
	}

	upkeep = upkeep.SetSelectedCategory(category)
	app.Upkeep = &upkeep

	return ViewSheets(app)
}

func HandleUpdate(app *App) (string, error) {
	if app.Params.Len() < 2 {
		return "", errors.New("invalid command, specify block id and category")
	}

	id, err := app.Params.GetInt(0)
	if err != nil {
		return "", err
	}

	timesheet := app.Timesheets[0].UpdateBlockCategory(id, app.Params.Get(1))
	app.Timesheets[0] = &timesheet

	return ViewSheets(app)
}

func HandleRestore(app *App) (string, error) {
	if app.Params.Len() == 0 {
		return "", errors.New("invalid command, specify block id")
	}

	id, err := app.Params.GetInt(0)
	if err != nil {
		return "", err
	}

	timesheet := app.Timesheets[0].RestoreBlock(id)
	app.Timesheets[0] = &timesheet

	return ViewSheets(app)
}

func HandleRemove(app *App) (string, error) {
	if app.Params.Len() == 0 {
		return "", errors.New("invalid command, specify block id")
	}

	id, err := app.Params.GetInt(0)
	if err != nil {
		return "", err
	}

	timesheet := app.Timesheets[0].RemoveBlock(id)
	app.Timesheets[0] = &timesheet

	return ViewSheets(app)
}

func HandleWrite(app *App) (string, error) {
	if app.Params.Len() < 2 {
		return "", errors.New("invalid command, specify category and duration")
	}

	cat := app.Params.Get(0)
	if app.Params.Get(1) == "fill" {
		quotum := app.Upkeep.GetTimesheetQuotum(*app.Timesheets[0])
		timesheet := app.Timesheets[0].Write(cat, quotum)
		app.Timesheets[0] = &timesheet
	} else {
		duration, err := time.ParseDuration(app.Params.Get(1))
		if err != nil {
			return "", err
		}

		timesheet := app.Timesheets[0].Write(cat, model.NewDuration().Set(duration))
		app.Timesheets[0] = &timesheet
	}

	return ViewSheets(app)
}

func HandleCategoryQuotum(app *App) (string, error) {
	if app.Params.Len() < 1 {
		return "", errors.New("invalid command, specify category and optional quotum")
	}

	cat := app.Params.Get(0)
	if app.Params.Len() == 1 {
		upkeep := app.Upkeep.SetCategoryMaxDayQuotum(cat, nil)
		app.Upkeep = &upkeep
	} else {
		d, err := time.ParseDuration(app.Params.Get(1))
		if err != nil {
			return "", err
		}
		upkeep := app.Upkeep.SetCategoryMaxDayQuotum(cat, &d)
		app.Upkeep = &upkeep
	}

	return ViewSheets(app)
}

func HandleQuotum(app *App) (string, error) {
	if app.Params.Len() == 0 {
		return "", errors.New("invalid command, specify weekday (0 = sunday) and optional quotum")
	}
	weekday, err := app.Params.GetInt(0)
	if err != nil {
		return "", err
	}
	if app.Params.Len() == 1 {
		upkeep := app.Upkeep.RemoveQuotumForDay(time.Weekday(weekday))
		app.Upkeep = &upkeep
		return fmt.Sprintf("removed quotum"), nil
	}

	duration, err := time.ParseDuration(app.Params.Get(1))
	if err != nil {
		return "", err
	}
	upkeep := app.Upkeep.SetQuotumForDay(time.Weekday(weekday), duration)
	app.Upkeep = &upkeep
	return fmt.Sprintf("updated quotum"), nil
}

func HandleVersion(app *App) (string, error) {
	return fmt.Sprintf("This is Upkeep version [%s]\n", app.Upkeep.Version), nil
}

func HandleFinalise(app *App) (string, error) {
	for i, t := range app.Timesheets {
		finalisedSheet := t.Finalise()
		app.Timesheets[i] = &finalisedSheet
	}

	return ViewDays(app)
}

func HandleUnfinalise(app *App) (string, error) {
	for i, t := range app.Timesheets {
		unfinalisedSheet := t.Unfinalise()
		app.Timesheets[i] = &unfinalisedSheet
	}

	return ViewDays(app)
}
