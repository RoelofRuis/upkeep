package app

import (
	"fmt"
	"regexp"
	"strconv"
	"upkeep/model"
)

var dateDefinition = regexp.MustCompile("^(-?[0-9]+)?([a-z]+)$")

// MakeDateRange shifts the given date based on the date selection definition given.
// It returns the shifted date and the number of selected days.
// A date definition can be one of the following:
// - An exact date given in the format yyyy-MM-DD
// - A string consisting of an optional integer number and a keyword/letter.
//   The keyword defines the base unit of time and selected duration.
//   The integer number defines how much to shift the base date through time.
//   Examples relative to an input date of today:
//     d  = today (length 1)
//   -1d  = yesterday (length 1)
//    1w  = next workweek starting monday (length 5)
//     wf = current workweek starting monday (length 7)
//   -3m  = three months ago starting first of month (length <number of days in that month>)
func MakeDateRange(baseDate model.Date, dateDef string) (model.Date, int, error) {
	shifts := 0
	numDays := 1

	matches := dateDefinition.FindStringSubmatch(dateDef)
	if len(matches) == 3 {
		if matches[1] != "" {
			i, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return baseDate, 0,  err
			}
			shifts = int(i)
		}
		dateDef = matches[2]
	}

	switch dateDef {
	case "day":
	case "d":
		baseDate = baseDate.ShiftDay(shifts)
		break

	case "week":
	case "w":
		baseDate = baseDate.PreviousMonday().ShiftDay(shifts * 7)
		numDays = 5
		break

	case "weekfull":
	case "wf":
		baseDate = baseDate.PreviousMonday().ShiftDay(shifts * 7)
		numDays = 7
		break

	case "month":
	case "m":
		baseDate = baseDate.FirstOfMonth().ShiftMonth(shifts)
		numDays = baseDate.DaysInMonth()
		break

	default:
		parsedDate, err := model.NewDateFromString(dateDef)
		if err != nil {
			return baseDate, 0, fmt.Errorf("invalid baseDate value '%s'", dateDef)
		}
		baseDate = parsedDate
		break
	}

	return baseDate, numDays, nil
}
