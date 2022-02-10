package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"
	"upkeep/model"
	"upkeep/model/repo"
)

type Repository repo.Repository

func HandleStart(args []string, editor *TimesheetEditor) (error, string) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Start(category)

	return nil, editor.View()
}

func HandleStop(args []string, editor *TimesheetEditor) (error, string) {
	editor.Stop()

	return nil, editor.View()
}

func HandleAbort(args []string, editor *TimesheetEditor) (error, string) {
	editor.Abort()

	return nil, editor.View()
}

func HandleSwitch(args []string, editor *TimesheetEditor) (error, string) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Switch(category)

	return nil, editor.View()
}

func HandleContinue(args []string, editor *TimesheetEditor) (error, string) {
	category := ""
	if len(args) > 0 {
		category = args[0]
	}

	editor.Continue(category)

	return nil, editor.View()
}

func HandleCategory(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) == 0 {
		return errors.New("no category specified"), ""
	}

	editor.Category(args[0])

	return nil, editor.View()
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

	return nil, editor.View()
}

func HandleDiscount(args []string, editor *TimesheetEditor) (error, string) {
	if len(args) < 2 {
		return errors.New("invalid command, specify category, type and optional arguments"), ""
	}

	cat := args[0]
	tpe := args[1]
	switch tpe {
	case "none":
		editor.RemoveDiscount(cat)
		break

	case "all":
		editor.AddDiscount(model.DiscountAll(cat))
		break

	case "max":
		if len(args) < 3 {
			return errors.New("specify max duration"), ""
		}
		d, err := time.ParseDuration(args[2])
		if err != nil {
			return err, ""
		}
		editor.AddDiscount(model.DiscountMax(cat, d))
		break

	default:
		return fmt.Errorf("unknown discount type '%s'", tpe), ""
	}

	return nil, editor.View()
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

func HandleVersion(args []string, editor *TimesheetEditor) (error, string) {
	return nil, fmt.Sprintf("Version: %s", editor.upkeep.Version)
}
