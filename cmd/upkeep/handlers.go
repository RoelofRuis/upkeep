package main

import (
	"errors"
	"fmt"
	"github.com/roelofruis/upkeep/internal/infra"
	"regexp"
	"time"

	"github.com/roelofruis/upkeep/internal/model"
)

func (a *App) Handle(f func(request *Request) (string, error)) infra.Handler {
	return func(params infra.Params) (string, error) {
		upkeep, err := a.Repo.Upkeep.Get()
		if err != nil {
			return "", err
		}

		date, numDays, err := MakeDateRange(model.NewDate(a.Clock.Now()), params)
		if err != nil {
			return "", err
		}

		dates := date.IterateNext(numDays)
		refSheets := make([]*model.Timesheet, len(dates))
		timesheets := make([]*model.Timesheet, len(dates))
		for i, day := range dates {
			sheet, err := a.Repo.Timesheets.GetForDate(day)
			if err != nil {
				return "", err
			}
			timesheets[i] = &sheet
			refSheets[i] = &sheet
		}

		req := &Request{
			Clock:      a.Clock,
			Params:     params,
			Upkeep:     &upkeep,
			Timesheets: timesheets,
		}

		s, err := f(req)
		if err != nil {
			return s, err
		}

		if req.Upkeep != &upkeep {
			if err := a.Repo.Upkeep.Insert(*req.Upkeep); err != nil {
				return "", err
			}
		}
		for i := 0; i < len(refSheets); i++ {
			if req.Timesheets[i] != refSheets[i] {
				if err := a.Repo.Timesheets.Insert(*req.Timesheets[i]); err != nil {
					return "", err
				}
			}
		}

		return s, nil
	}
}

func HandleStart(req *Request) (string, error) {
	category := req.Params.Get(0)

	_, err := HandleStop(req)
	if err != nil {
		return "", err
	}

	now := req.Clock.Now()
	sheet := req.Timesheets[0].Start(now)

	quotum := req.Upkeep.GetWeekdayQuotum(now.Weekday())
	sheet = sheet.SetQuotum(quotum)

	req.Timesheets[0] = &sheet

	if category != "" {
		_, err := HandleSet(req)
		if err != nil {
			return "", err
		}
	}

	return ViewSheets(req)
}

func HandleStop(req *Request) (string, error) {
	sheet := req.Timesheets[0].Stop(req.Clock.Now(), req.Upkeep.GetSelectedCategory().Name)
	req.Timesheets[0] = &sheet

	return ViewSheets(req)
}

func HandleAbort(req *Request) (string, error) {
	sheet := req.Timesheets[0].Abort()
	req.Timesheets[0] = &sheet

	return ViewSheets(req)
}

func HandleSwitch(req *Request) (string, error) {
	_, err := HandleStop(req)
	if err != nil {
		return "", err
	}

	upkeep := req.Upkeep.ShiftSelectedCategory()
	req.Upkeep = &upkeep

	return HandleStart(req)
}

func HandleContinue(req *Request) (string, error) {
	_, err := HandleStop(req)
	if err != nil {
		return "", err
	}

	upkeep := req.Upkeep.UnshiftSelectedCategory()
	req.Upkeep = &upkeep
	return HandleStart(req)
}

var validCategory = regexp.MustCompile(`^[a-zA-Z0-9_-]+(\.[a-zA-Z0-9_-]+)?$`)

func HandleSet(req *Request) (string, error) {
	if req.Params.Len() == 0 {
		return "", errors.New("no category specified")
	}

	upkeep := *req.Upkeep
	category := req.Params.Get(0)

	if !validCategory.MatchString(category) {
		return ViewSheets(req)
	}

	upkeep = upkeep.SetSelectedCategory(category)
	req.Upkeep = &upkeep

	return ViewSheets(req)
}

func HandleUpdate(req *Request) (string, error) {
	if req.Params.Len() < 2 {
		return "", errors.New("invalid command, specify block id and category")
	}

	id, err := req.Params.GetInt(0)
	if err != nil {
		return "", err
	}

	timesheet := req.Timesheets[0].UpdateBlockCategory(id, req.Params.Get(1))
	req.Timesheets[0] = &timesheet

	return ViewSheets(req)
}

func HandleRestore(req *Request) (string, error) {
	if req.Params.Len() == 0 {
		return "", errors.New("invalid command, specify block id")
	}

	id, err := req.Params.GetInt(0)
	if err != nil {
		return "", err
	}

	timesheet := req.Timesheets[0].RestoreBlock(id)
	req.Timesheets[0] = &timesheet

	return ViewSheets(req)
}

func HandleRemove(req *Request) (string, error) {
	if req.Params.Len() == 0 {
		return "", errors.New("invalid command, specify block id")
	}

	id, err := req.Params.GetInt(0)
	if err != nil {
		return "", err
	}

	timesheet := req.Timesheets[0].RemoveBlock(id)
	req.Timesheets[0] = &timesheet

	return ViewSheets(req)
}

func HandleWrite(req *Request) (string, error) {
	if req.Params.Len() < 2 {
		return "", errors.New("invalid command, specify category and duration")
	}

	cat := req.Params.Get(0)
	if req.Params.Get(1) == "fill" {
		quotum := req.Upkeep.GetTimesheetQuotum(*req.Timesheets[0])
		timesheet := req.Timesheets[0].Write(cat, quotum)
		req.Timesheets[0] = &timesheet
	} else {
		duration, err := time.ParseDuration(req.Params.Get(1))
		if err != nil {
			return "", err
		}

		timesheet := req.Timesheets[0].Write(cat, model.NewDuration().Set(duration))
		req.Timesheets[0] = &timesheet
	}

	return ViewSheets(req)
}

func HandleCategoryQuotum(req *Request) (string, error) {
	if req.Params.Len() < 1 {
		return "", errors.New("invalid command, specify category and optional quotum")
	}

	cat := req.Params.Get(0)
	if req.Params.Len() == 1 {
		upkeep := req.Upkeep.SetCategoryMaxDayQuotum(cat, nil)
		req.Upkeep = &upkeep
	} else {
		d, err := time.ParseDuration(req.Params.Get(1))
		if err != nil {
			return "", err
		}
		upkeep := req.Upkeep.SetCategoryMaxDayQuotum(cat, &d)
		req.Upkeep = &upkeep
	}

	return ViewSheets(req)
}

func HandleQuotum(req *Request) (string, error) {
	if req.Params.Len() == 0 {
		return "", errors.New("invalid command, specify weekday (0 = sunday) and optional quotum")
	}
	weekday, err := req.Params.GetInt(0)
	if err != nil {
		return "", err
	}
	if req.Params.Len() == 1 {
		upkeep := req.Upkeep.RemoveQuotumForDay(time.Weekday(weekday))
		req.Upkeep = &upkeep
		return fmt.Sprintf("removed quotum"), nil
	}

	duration, err := time.ParseDuration(req.Params.Get(1))
	if err != nil {
		return "", err
	}
	upkeep := req.Upkeep.SetQuotumForDay(time.Weekday(weekday), duration)
	req.Upkeep = &upkeep
	return fmt.Sprintf("updated quotum"), nil
}

func HandleVersion(req *Request) (string, error) {
	return fmt.Sprintf("This is Upkeep version [%s]\n", req.Upkeep.Version), nil
}

func HandleFinalise(req *Request) (string, error) {
	for i, t := range req.Timesheets {
		finalisedSheet := t.Finalise()
		req.Timesheets[i] = &finalisedSheet
	}

	return ViewDays(req)
}

func HandleUnfinalise(req *Request) (string, error) {
	for i, t := range req.Timesheets {
		unfinalisedSheet := t.Unfinalise()
		req.Timesheets[i] = &unfinalisedSheet
	}

	return ViewDays(req)
}
