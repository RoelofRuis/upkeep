package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"upkeep/model/repo"
)

type Repository repo.Repository

func HandleStart(args []string, editor *TimesheetEditor) (error, string) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Start(category)

	return nil, editor.Day()
}

func HandleStop(args []string, editor *TimesheetEditor) (error, string) {
	editor.Stop()

	return nil, editor.Day()
}

func HandleAbort(args []string, editor *TimesheetEditor) (error, string) {
	editor.Abort()

	return nil, editor.Day()
}

func HandleSwitch(args []string, editor *TimesheetEditor) (error, string) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Switch(category)

	return nil, editor.Day()
}

func HandleContinue(args []string, editor *TimesheetEditor) (error, string) {
	editor.Continue()

	return nil, editor.Day()
}

func HandleCategory(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no category specified"), ""
	}

	editor.Category(args[0])

	return nil, editor.Day()
}

func HandleDay(args []string, editor *TimesheetEditor) (error, string) {
	return nil, editor.Day()
}

func HandleRemove(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no id given"), ""
	}

	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err, ""
	}

	editor.Remove(int(id))

	return nil, editor.Day()
}

func HandleExclude(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no category given"), ""
	}

	editor.Exclude(args[0])

	return nil, editor.Day()
}

func HandleInclude(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no category given"), ""
	}

	editor.Inlcude(args[0])

	return nil, editor.Day()
}

func HandleQuotum(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("too few arguments"), ""
	}
	weekday, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err, ""
	}
	if len(args) == 1 {
		editor.AdjustQuotum(time.Weekday(weekday), nil)
		return nil, fmt.Sprintf("removed quotum")
	}

	duration, err := time.ParseDuration(args[1])
	if err != nil {
		return err, ""
	}
	editor.AdjustQuotum(time.Weekday(weekday), &duration)
	return nil, fmt.Sprintf("updated quotum")
}
