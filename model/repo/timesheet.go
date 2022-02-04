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
	CreatedAt string      `json:"created_at"`
	NextId    int         `json:"next_id"`
	LastStart string      `json:"last_start"`
	Blocks    []blockJson `json:"blocks"`
}

type blockJson struct {
	Id    int    `json:"id"`
	Start string `json:"start"`
	End   string `json:"end"`
	Tags  string `json:"tags"`
}

func (r *TimesheetRepository) GetForDay(t time.Time) (*model.Timesheet, error) {
	input := timesheetJson{
		CreatedAt: t.Format(model.LayoutDateHour),
		NextId:    0,
		LastStart: "",
		Blocks:    nil,
	}

	if err := r.FileIO.Read(fmt.Sprintf("/sheet/%s.json", t.Format(model.LayoutHour)), &input); err != nil {
		return nil, err
	}

	createdTime, err := time.Parse(model.LayoutDateHour, input.CreatedAt)
	if err != nil {
		return nil, err
	}

	sheet := model.NewTimesheet(createdTime)
	sheet.NextId = input.NextId

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
			Id:    blockData.Id,
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
			Id:    block.Id,
			Start: block.Start.Format(model.LayoutDateHour),
			End:   block.End.Format(model.LayoutDateHour),
			Tags:  block.Tags.String(),
		})
	}

	output := timesheetJson{
		CreatedAt: m.Created.Format(model.LayoutDateHour),
		NextId:    m.NextId,
		LastStart: m.LastStart.Format(model.LayoutDateHour),
		Blocks:    blocks,
	}

	if err := r.FileIO.Write(fmt.Sprintf("/sheet/%s.json", m.Created.Format(model.LayoutDate)), output); err != nil {
		return err
	}

	return nil
}