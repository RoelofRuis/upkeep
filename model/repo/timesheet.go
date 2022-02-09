package repo

import (
	"fmt"
	"time"
	"upkeep/infra"
	"upkeep/model"
)

type TimesheetRepository struct {
	FileIO infra.FileIO
}

type timesheetJson struct {
	NextId    int         `json:"next_id"`
	LastStart string      `json:"last_start"`
	Blocks    []blockJson `json:"blocks"`
	Quotum    string      `json:"quotum"`
}

type blockJson struct {
	Id       int    `json:"id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Category string `json:"category"`
}

func (r *TimesheetRepository) GetForDate(date model.Date) (model.Timesheet, error) {
	input := timesheetJson{
		NextId:    0,
		LastStart: "",
		Blocks:    nil,
		Quotum:    "",
	}

	if err := r.FileIO.Read(fmt.Sprintf("/sheet/%s.json", date), &input); err != nil {
		return model.Timesheet{}, err
	}

	sheet := model.NewTimesheet(date)
	sheet.NextId = input.NextId

	if input.Quotum != "" {
		quotum, err := time.ParseDuration(input.Quotum)
		if err != nil {
			return model.Timesheet{}, err
		}
		sheet.Quotum = quotum
	}

	lastStart, err := model.NewMomentFromString(input.LastStart)
	if err != nil {
		return model.Timesheet{}, err
	}
	sheet.LastStart = lastStart

	var blocks []model.TimeBlock
	for _, blockData := range input.Blocks {
		start, err := model.NewMomentFromString(blockData.Start)
		if err != nil {
			return model.Timesheet{}, err
		}
		end, err := model.NewMomentFromString(blockData.End)
		if err != nil {
			return model.Timesheet{}, err
		}
		block := model.TimeBlock{
			Id:       blockData.Id,
			Start:    start,
			End:      end,
			Category: blockData.Category,
		}
		blocks = append(blocks, block)
	}

	sheet.Blocks = blocks

	return sheet, nil
}

func (r *TimesheetRepository) Insert(m model.Timesheet) error {
	var blocks []blockJson

	for _, block := range m.Blocks {
		blocks = append(blocks, blockJson{
			Id:       block.Id,
			Start:    block.Start.Format(model.LayoutDateHour),
			End:      block.End.Format(model.LayoutDateHour),
			Category: block.Category,
		})
	}

	output := timesheetJson{
		NextId:    m.NextId,
		LastStart: m.LastStart.Format(model.LayoutDateHour),
		Blocks:    blocks,
		Quotum:    m.Quotum.String(),
	}

	if err := r.FileIO.Write(fmt.Sprintf("/sheet/%s.json", m.Date), output); err != nil {
		return err
	}

	return nil
}
