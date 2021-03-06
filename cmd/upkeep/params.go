package main

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

var dateDefinition = regexp.MustCompile("^(-?[0-9]+)?([a-zA-Z]+)$")

// GroupCategories determines whether to group categories or use the full category names
// Specify g:g if you want to group by category group
// specify g:c if you want to group by individual category
func GroupCategories(params infra.Params) bool {
	return params.GetNamed("g", "c") == "g"
}

// MakeDateRange shifts the given date based on the date parameter provided.
// It returns the shifted date and the number of selected days.
// A date definition can be one of the following:
// - An exact date given in the format yyyy-MM-DD
// - A string consisting of an optional integer number and a keyword/letter.
//   The keyword defines the base unit of time and selected duration.
//   The integer number defines how much to shift the base date through time.
//   Examples relative to an input date of today:
//     d = today (length 1)
//   -1d = yesterday (length 1)
//    1w = next workweek starting monday (length 7)
//     w = current workweek starting monday (length 7)
//    wr = rolling week: past seven days, ending today (length 7)
//   -3m = three months ago starting first of month (length <number of days in that month>)
func MakeDateRange(baseDate model.Date, params infra.Params) (model.Date, int, error) {
	dateDef := params.GetNamed("d", "day")

	shifts := 0
	numDays := 1

	matches := dateDefinition.FindStringSubmatch(dateDef)
	if len(matches) == 3 {
		if matches[1] != "" {
			i, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				return baseDate, 0, err
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
		numDays = 7
		break

	case "weekrolling":
	case "wr":
		baseDate = baseDate.ShiftDay((shifts-1)*7 + 1)
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
			return baseDate, 0, fmt.Errorf("invalid date value '%s'", dateDef)
		}
		baseDate = parsedDate
		break
	}

	return baseDate, numDays, nil
}
