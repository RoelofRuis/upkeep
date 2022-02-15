package app

import (
	"errors"
	"fmt"
	"time"
	"upkeep/infra"
	"upkeep/model"
	"upkeep/model/repo"
)

type Repository repo.Repository

func HandleStart(params infra.Params, editor *TimesheetEditor) (string, error) {
	editor.Start(params.Get(0))

	return editor.View(), nil
}

func HandleStop(params infra.Params, editor *TimesheetEditor) (string, error) {
	editor.Stop()

	return editor.View(), nil
}

func HandleAbort(params infra.Params, editor *TimesheetEditor) (string, error) {
	editor.Abort()

	return editor.View(), nil
}

func HandleSwitch(params infra.Params, editor *TimesheetEditor) (string, error) {
	editor.Switch(params.Get(0))

	return editor.View(), nil
}

func HandleContinue(params infra.Params, editor *TimesheetEditor) (string, error) {
	editor.Continue(params.Get(0))

	return editor.View(), nil
}

func HandleSet(params infra.Params, editor *TimesheetEditor) (string, error) {
	if params.Len() == 0 {
		return "", errors.New("no category specified")
	}

	editor.Category(params.Get(0))

	return editor.View(), nil
}

func HandleUpdate(params infra.Params, editor *TimesheetEditor) (string, error) {
	if params.Len() < 2 {
		return "", errors.New("invalid command, specify block id and category")
	}

	id, err := params.GetInt(0)
	if err != nil {
		return "", err
	}

	editor.Update(id, params.Get(1))

	return editor.View(), nil
}

func HandleRemove(params infra.Params, editor *TimesheetEditor) (string, error) {
	if params.Len() == 0 {
		return "", errors.New("invalid command, specify block id")
	}

	id, err := params.GetInt(0)
	if err != nil {
		return "", err
	}

	editor.Remove(id)

	return editor.View(), nil
}

func HandleWrite(params infra.Params, editor *TimesheetEditor) (string, error) {
	if params.Len() < 2 {
		return "", errors.New("invalid command, specify category and duration")
	}

	cat := params.Get(0)
	if params.Get(1) == "fill" {
		quotum := editor.upkeep.GetTimesheetQuotum(*editor.timesheet)
		editor.Write(cat, quotum)
	} else {
		duration, err := time.ParseDuration(params.Get(1))
		if err != nil {
			return "", err
		}

		editor.Write(cat, model.NewDuration().Set(duration))
	}

	return editor.View(), nil
}

func HandleCategoryQuotum(params infra.Params, editor *TimesheetEditor) (string, error) {
	if params.Len() < 1 {
		return "", errors.New("invalid command, specify category and optional quotum")
	}

	cat := params.Get(0)
	if params.Len() == 1 {
		editor.SetCategoryMaxDayQuotum(cat, nil)
	} else {
		d, err := time.ParseDuration(params.Get(1))
		if err != nil {
			return "", err
		}
		editor.SetCategoryMaxDayQuotum(cat, &d)
	}

	return editor.View(), nil
}

func HandleQuotum(params infra.Params, editor *TimesheetEditor) (string, error) {
	if params.Len() == 0 {
		return "", errors.New("invalid command, specify weekday (0 = sunday) and optional quotum")
	}
	weekday, err := params.GetInt(0)
	if err != nil {
		return "", err
	}
	if params.Len() == 1 {
		editor.AdjustQuotum(time.Weekday(weekday), nil)
		return fmt.Sprintf("removed quotum"), nil
	}

	duration, err := time.ParseDuration(params.Get(1))
	if err != nil {
		return "", err
	}
	editor.AdjustQuotum(time.Weekday(weekday), &duration)
	return fmt.Sprintf("updated quotum"), nil
}

func HandleVersion(params infra.Params, editor *TimesheetEditor) (string, error) {
	return fmt.Sprintf("Version: %s", editor.upkeep.Version), nil
}
