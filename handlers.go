package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func (app *application) handleStart(args []string) (error, string) {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err, ""
	}

	timesheet.Start(time.Now())

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err, ""
	}

	return nil, "started"
}

func (app *application) handleStop(args []string) (error, string) {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err, ""
	}

	timesheet.Stop(time.Now())

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err, ""
	}

	return nil, "stopped"
}

var validTag = regexp.MustCompile(`^[+-]?[a-z]*$`)

func (app *application) handleTag(args []string) (error, string) {
	if len(args) < 1 {
		return errors.New("no tag specified"), ""
	}

	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err, ""
	}

	for _, tag := range args {
		if !validTag.MatchString(tag) {
			return fmt.Errorf("invalid tag '%s'", tag), ""
		}
		if strings.HasPrefix(tag, "-") {
			timesheet.DetachTag(strings.TrimPrefix(tag, "-"))
		} else {
			timesheet.AttachTag(strings.TrimPrefix(tag, "+"))
		}
	}

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err, ""
	}

	return nil, "tags updated"
}

func (app *application) handleShow(args []string) (error, string) {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err, ""
	}

	msg := fmt.Sprintf("Day: %s", timesheet.Day)
	return nil, msg
}
