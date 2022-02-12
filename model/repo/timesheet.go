package repo

import (
	"fmt"
	"upkeep/infra"
	"upkeep/model"
)

type TimesheetRepository struct {
	FileIO infra.FileIO
}

type timesheetJson struct {
	NextId    int                `json:"next_id"`
	LastStart model.Moment       `json:"last_start"`
	Blocks    []blockJson        `json:"blocks"`
	Quotum    infra.JSONDuration `json:"quotum"`
}

type blockJson struct {
	Id       int          `json:"id"`
	Start    model.Moment `json:"start"`
	End      model.Moment `json:"end"`
	Category string       `json:"category"`
}

func (r *TimesheetRepository) GetForDate(date model.Date) (model.Timesheet, error) {
	input := timesheetJson{
		NextId:    0,
		LastStart: model.NewMoment(),
		Blocks:    nil,
		Quotum:    infra.NewDuration(),
	}

	if err := r.FileIO.Read(fmt.Sprintf("/sheet/%s.json", date), &input); err != nil {
		return model.Timesheet{}, err
	}

	sheet := model.NewTimesheet(date)
	sheet.NextId = input.NextId
	sheet.Quotum = input.Quotum.Unpack()
	sheet.LastStart = input.LastStart

	var blocks []model.TimeBlock
	for _, blockData := range input.Blocks {
		blocks = append(blocks, model.TimeBlock{
			Id:       blockData.Id,
			Start:    blockData.Start,
			End:      blockData.End,
			Category: blockData.Category,
		})
	}

	sheet.Blocks = blocks

	return sheet, nil
}

func (r *TimesheetRepository) Insert(m model.Timesheet) error {
	var blocks []blockJson

	for _, block := range m.Blocks {
		blocks = append(blocks, blockJson{
			Id:       block.Id,
			Start:    block.Start,
			End:      block.End,
			Category: block.Category,
		})
	}

	output := timesheetJson{
		NextId:    m.NextId,
		LastStart: m.LastStart,
		Blocks:    blocks,
		Quotum:    infra.JSONDuration(m.Quotum),
	}

	if err := r.FileIO.Write(fmt.Sprintf("/sheet/%s.json", m.Date), output); err != nil {
		return err
	}

	return nil
}
