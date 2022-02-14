package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"upkeep/model/repo"
)

type Repository repo.Repository

func HandleStart(args []string, editor *TimesheetEditor) (string, error) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Start(category)

	return editor.View(), nil
}

func HandleStop(args []string, editor *TimesheetEditor) (string, error) {
	editor.Stop()

	return editor.View(), nil
}

func HandleAbort(args []string, editor *TimesheetEditor) (string, error) {
	editor.Abort()

	return editor.View(), nil
}

func HandleSwitch(args []string, editor *TimesheetEditor) (string, error) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Switch(category)

	return editor.View(), nil
}

func HandleContinue(args []string, editor *TimesheetEditor) (string, error) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Continue(category)

	return editor.View(), nil
}

func HandleSet(args []string, editor *TimesheetEditor) (string, error) {
	if len(args) == 0 {
		return "", errors.New("no category specified")
	}

	editor.Category(args[0])

	return editor.View(), nil
}

func HandleUpdate(args []string, editor *TimesheetEditor) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid command, specify block id and category")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "", err
	}

	editor.Update(int(id), args[1])

	return editor.View(), nil
}

func HandleRemove(args []string, editor *TimesheetEditor) (string, error) {
	if len(args) == 0 {
		return "", errors.New("invalid command, specify block id")
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "", err
	}

	editor.Remove(int(id))

	return editor.View(), nil
}

func HandleWrite(args []string, editor *TimesheetEditor) (string, error) {
	if len(args) < 2 {
		return "", errors.New("invalid command, specify category and duration")
	}

	cat := args[0]
	duration, err := time.ParseDuration(args[1])
	if err != nil {
		return "", err
	}

	editor.Write(cat, duration)

	return editor.View(), nil
}

func HandleCategoryQuotum(args []string, editor *TimesheetEditor) (string, error) {
	if len(args) < 1 {
		return "", errors.New("invalid command, specify category and optional quotum")
	}

	cat := args[0]
	if len(args) == 1 {
		editor.SetCategoryMaxDayQuotum(cat, nil)
	} else {
		d, err := time.ParseDuration(args[1])
		if err != nil {
			return "", err
		}
		editor.SetCategoryMaxDayQuotum(cat, &d)
	}

	return editor.View(), nil
}

func HandleQuotum(args []string, editor *TimesheetEditor) (string, error) {
	if len(args) == 0 {
		return "", errors.New("invalid command, specify weekday (0 = sunday) and optional quotum")
	}
	weekday, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return "", err
	}
	if len(args) == 1 {
		editor.AdjustQuotum(time.Weekday(weekday), nil)
		return fmt.Sprintf("removed quotum"), nil
	}

	duration, err := time.ParseDuration(args[1])
	if err != nil {
		return "", err
	}
	editor.AdjustQuotum(time.Weekday(weekday), &duration)
	return fmt.Sprintf("updated quotum"), nil
}

func HandleVersion(args []string, editor *TimesheetEditor) (string, error) {
	return fmt.Sprintf("Version: %s", editor.upkeep.Version), nil
}
