package main

import (
	"errors"
	"fmt"
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

	err = timesheet.Start(time.Now())
	if err != nil {
		return err
	}

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

	err = timesheet.Stop(time.Now())
	if err != nil {
		return err
	}

	err = app.timesheetRepository.Insert(timesheet)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) handleTag(args []string) error {
	if len(args) < 1 {
		return errors.New("no tag specified")
	}

	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err
	}

	// TODO: process every tag, parse them and attach or detach

	err = timesheet.TagActiveBlock(args[0])
	if err != nil {
		return err
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
