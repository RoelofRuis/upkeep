package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"timesheet/model"
)

type domain struct {
	upkeep    *model.Upkeep
	timesheet *model.Timesheet
}

func (app *application) withDomain(f func(args []string, domain domain) (error, string)) func(args []string) (error, string) {
	return func(args []string) (error, string) {
		upkeep, err := app.upkeepRepository.Get()
		if err != nil {
			return err, ""

		}
		timesheet, err := app.timesheetRepository.GetForDay(time.Now())
		if err != nil {
			return err, ""
		}

		err, s := f(args, domain{upkeep: upkeep, timesheet: timesheet})
		if err != nil {
			return err, s
		}

		if err := app.upkeepRepository.Insert(upkeep); err != nil {
			return err, ""
		}
		if err := app.timesheetRepository.Insert(timesheet); err != nil {
			return err, ""
		}

		return nil, s
	}
}

func (app *application) handlePurge(args []string) (error, string) {
	timesheet, err := app.timesheetRepository.GetForDay(time.Now())
	if err != nil {
		return err, ""
	}

	if err := app.timesheetRepository.Delete(timesheet); err != nil {
		return err, ""
	}

	return nil, "purged"
}

func handleStart(args []string, domain domain) (error, string) {
	domain.timesheet.Start(time.Now())

	return nil, "started new block"
}

func handleStop(args []string, domain domain) (error, string) {
	domain.timesheet.Stop(time.Now(), domain.upkeep.GetTags())

	return nil, "stopped active block"
}

func handleSwitch(args []string, domain domain) (error, string) {
	domain.timesheet.Stop(time.Now(), domain.upkeep.GetTags())

	if len(args) == 1 && args[0] == "-" {
		domain.upkeep.UnshiftTags()
	} else {
		domain.upkeep.ShiftTags()
	}

	domain.timesheet.Start(time.Now())

	return nil, "switched"
}

var validTag = regexp.MustCompile(`^[+-]?[a-z]*$`)

func handleTag(args []string, domain domain) (error, string) {
	if len(args) < 1 {
		return errors.New("no tag specified"), ""
	}

	for _, tag := range args {
		if !validTag.MatchString(tag) {
			return fmt.Errorf("invalid tag '%s'", tag), ""
		}
		if strings.HasPrefix(tag, "-") {
			domain.upkeep.RemoveTag(strings.TrimPrefix(tag, "-"))
		} else {
			domain.upkeep.AddTag(strings.TrimPrefix(tag, "+"))
		}
	}

	return nil, "tags updated"
}

func handleShow(args []string, domain domain) (error, string) {
	var lines []string
	lines = append(lines, fmt.Sprintf("> %s [%s]", domain.timesheet.Day, domain.upkeep.Tags.String()))
	for _, block := range domain.timesheet.Blocks {
		blockString := fmt.Sprintf("%s - %s [%s]", block.Start.HourString(), block.End.HourString(), block.Tags.String())
		lines = append(lines, blockString)
	}
	if domain.timesheet.LastStart.IsStarted() {
		activeBlockString := fmt.Sprintf("%s -   ?   [%s]", domain.timesheet.LastStart.HourString(), domain.upkeep.GetTags().String())
		lines = append(lines, activeBlockString)
	}
	return nil, strings.Join(lines, "\n")
}
