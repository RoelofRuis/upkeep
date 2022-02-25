package repo

import (
	"fmt"

	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model"
)

type TimesheetRepository struct {
	FileIO infra.FileIO
}

type timesheetJson struct {
	Date      model.Date     `json:"date"`
	NextId    int            `json:"next_id"`
	LastStart model.Moment   `json:"last_start"`
	Blocks    []blockJson    `json:"blocks"`
	Quotum    model.Duration `json:"quotum"`
	Finalised bool           `json:"finalised"`
}

type blockJson struct {
	Id       int            `json:"id"`
	Category string         `json:"category"`
	Deleted  bool           `json:"deleted"`
	Type     string         `json:"type"`
	Start    model.Moment   `json:"start"`
	End      model.Moment   `json:"end"`
	Duration model.Duration `json:"duration"`
}

func (r *TimesheetRepository) GetForDate(date model.Date) (model.Timesheet, error) {
	input := timesheetJson{
		NextId:    0,
		LastStart: model.NewMoment(),
		Blocks:    nil,
		Quotum:    model.NewDuration(),
		Finalised: false,
	}

	if err := r.FileIO.Read(filenameForDate(date), &input); err != nil {
		return model.Timesheet{}, err
	}

	sheet := model.NewTimesheet(date)
	sheet.NextId = input.NextId
	sheet.Quotum = input.Quotum
	sheet.LastStart = input.LastStart
	sheet.Finalised = input.Finalised

	var blocks []model.TimeBlock
	for _, blockData := range input.Blocks {
		switch blockData.Type {
		case model.TypeDuration:
			blocks = append(blocks, model.NewBlockWithDuration(
				blockData.Id,
				blockData.Category,
				blockData.Deleted,
				blockData.Duration,
			))
		default:
			blocks = append(blocks, model.NewBlockWithTime(
				blockData.Id,
				blockData.Category,
				blockData.Deleted,
				blockData.Start,
				blockData.End,
			))
		}

	}

	sheet.Blocks = blocks

	return sheet, nil
}

func (r *TimesheetRepository) Insert(m model.Timesheet) error {
	var blocks []blockJson

	for _, block := range m.Blocks {
		blocks = append(blocks, blockJson{
			Id:       block.Id,
			Type:     block.Type,
			Start:    block.WithTime.Start,
			End:      block.WithTime.End,
			Duration: block.WithDuration.Duration,
			Category: block.Category,
			Deleted:  block.Deleted,
		})
	}

	output := timesheetJson{
		Date:      m.Date,
		NextId:    m.NextId,
		LastStart: m.LastStart,
		Blocks:    blocks,
		Quotum:    m.Quotum,
		Finalised: m.Finalised,
	}

	if err := r.FileIO.Write(filenameForDate(m.Date), output); err != nil {
		return err
	}

	return nil
}

func filenameForDate(date model.Date) string {
	return fmt.Sprintf(fmt.Sprintf("/sheet/%s.json", date))
}
