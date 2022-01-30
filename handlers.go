package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

func (app *application) handle(args []string) {
	if len(args) == 0 {
		fmt.Printf("timesheet command\n - start\n - stop\n - tag\n - show\n")
		return
	}

	var handlerError error
	switch args[0] {
	case "start":
		handlerError = app.handleStart(args[1:])
	case "stop":
		handlerError = app.handleStop(args[1:])
	case "tag":
		handlerError = app.handleTag(args[1:])
	case "show":
		handlerError = app.handleShow(args[1:])
	default:
		handlerError = fmt.Errorf("unknown command '%s'", args[0])
	}

	if handlerError != nil {
		fmt.Printf("error: %s\n", handlerError.Error())
	}
}

func (app *application) handleStart(args []string) error {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err
	}

	timesheet.Start(time.Now())

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) handleStop(args []string) error {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err
	}

	timesheet.Stop(time.Now())

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err
	}

	return nil
}

var validTag = regexp.MustCompile(`^[+-]?[a-z]*$`)

func (app *application) handleTag(args []string) error {
	if len(args) < 1 {
		return errors.New("no tag specified")
	}

	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err
	}

	for _, tag := range args {
		if !validTag.MatchString(tag) {
			return fmt.Errorf("invalid tag '%s'", tag)
		}
		if strings.HasPrefix(tag, "-") {
			timesheet.UntagActiveBlock(strings.TrimPrefix(tag, "-"))
		} else {
			timesheet.TagActiveBlock(strings.TrimPrefix(tag, "+"))
		}
	}

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) handleShow(args []string) error {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err
	}

	PrettyPrint(timesheet)
	return nil
}
