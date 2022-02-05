package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"timesheet/model/repo"
)

type Repository repo.Repository

func HandlePurge(args []string, editor *TimesheetEditor) (error, string) {
	editor.Purge()

	return nil, editor.Show()
}

func HandleStart(args []string, editor *TimesheetEditor) (error, string) {
	editor.Start(args)

	return nil, editor.Show()
}

func HandleStop(args []string, editor *TimesheetEditor) (error, string) {
	editor.Stop()

	return nil, editor.Show()
}

func HandleAbort(args []string, editor *TimesheetEditor) (error, string) {
	editor.Abort()

	return nil, editor.Show()
}

func HandleSwitch(args []string, editor *TimesheetEditor) (error, string) {
	editor.Switch(args)

	return nil, editor.Show()
}

func HandleTag(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no tags specified"), ""
	}

	editor.Tag(args)

	return nil, editor.Show()
}

func HandleShow(args []string, editor *TimesheetEditor) (error, string) {
	return nil, editor.Show()
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
