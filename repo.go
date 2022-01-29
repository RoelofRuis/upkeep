package main

import (
	"encoding/csv"
	"os"
)

type TimesheetRepository struct {
}

func (r *TimesheetRepository) Insert(t *Timesheet) error {
	f, err := os.Create("sheet.csv")
	if err != nil {
		return err
	}
	defer f.Close()

	csvWriter := csv.NewWriter(f)
	var rows [][]string
	for _, block := range t.Blocks {
		rows = append(rows, []string{block.Start.String(), block.End.String()})
	}

	err = csvWriter.WriteAll(rows)
	if err != nil {
		return err
	}

	return nil
}

func (r *TimesheetRepository) Get() (*Timesheet, error) {
	// TODO: timesheet per day
	f, err := os.Open("sheet.csv")
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return &Timesheet{}, nil
		default:
			return nil, err
		}
	}

	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	var blocks []TimeBlock

	for _, row := range rows {
		start, err := NewMomentFromString(row[0])
		if err != nil {
			return nil, err
		}
		end, err := NewMomentFromString(row[1])
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, TimeBlock{
			Start: start,
			End:   end,
		})
	}

	sheet := &Timesheet{
		Blocks: blocks,
	}

	return sheet, nil
}
