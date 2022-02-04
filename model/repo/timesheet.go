package repo

import (
	"fmt"
	"time"
	"timesheet/infra"
	"timesheet/model"
)

type TimesheetRepository struct {
	FileIO infra.FileIO
}

type timesheetJson struct {
	Day       string      `json:"day"`
	NextId    int         `json:"next_id"`
	LastStart string      `json:"last_start"`
	Blocks    []blockJson `json:"blocks"`
}

type blockJson struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Tags  string `json:"tags"`
}

func (r *TimesheetRepository) GetForDay(t time.Time) (*model.Timesheet, error) {
	day := t.Format("2006-01-02")

	input := timesheetJson{
		Day:       day,
		NextId:    0,
		LastStart: "",
		Blocks:    nil,
	}

	if err := r.FileIO.Read(fmt.Sprintf("/sheet/%s.json", day), &input); err != nil {
		return nil, err
	}

	sheet := model.NewTimesheet(input.Day)

	lastStart, err := model.NewMomentFromString(input.LastStart)
	if err != nil {
		return nil, err
	}
	sheet.LastStart = lastStart

	var blocks []model.TimeBlock
	for _, blockData := range input.Blocks {
		start, err := model.NewMomentFromString(blockData.Start)
		if err != nil {
			return nil, err
		}
		end, err := model.NewMomentFromString(blockData.End)
		if err != nil {
			return nil, err
		}
		block := model.TimeBlock{
			Start: start,
			End:   end,
			Tags:  model.NewTagSetFromString(blockData.Tags),
		}
		blocks = append(blocks, block)
	}

	sheet.Blocks = blocks

	return sheet, nil
}

func (r *TimesheetRepository) Insert(m *model.Timesheet) error {
	var blocks []blockJson

	for _, block := range m.Blocks {
		blocks = append(blocks, blockJson{
			Start: block.Start.String(),
			End:   block.End.String(),
			Tags:  block.Tags.String(),
		})
	}

	output := timesheetJson{
		Day:       m.Day,
		NextId:    m.NextId,
		LastStart: m.LastStart.String(),
		Blocks:    blocks,
	}

	if err := r.FileIO.Write(fmt.Sprintf("/sheet/%s.json", m.Day), output); err != nil {
		return err
	}

	return nil
}

func (r *TimesheetRepository) Delete(m *model.Timesheet) error {
	return r.FileIO.Delete(fmt.Sprintf("/sheet/%s.json", m.Day))
}
