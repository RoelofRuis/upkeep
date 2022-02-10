package repo

import (
	"time"
	"upkeep/infra"
	"upkeep/model"
)

type UpkeepRepository struct {
	FileIO infra.FileIO
}

type upkeepJson struct {
	Version        string         `json:"version"`
	ActiveCategory string         `json:"active_category"`
	Quota          map[int]string `json:"quota"`
	Discounts      []discountJson `json:"discounts"`
}

type discountJson struct {
	Category string `json:"category"`
	Argument string `json:"argument"`
}

const VERSION = "0.2"

func (r *UpkeepRepository) Get() (model.Upkeep, error) {
	input := upkeepJson{
		Version: "0.1",
	}

	if err := r.FileIO.Read("upkeep.json", &input); err != nil {
		return model.Upkeep{}, err
	}

	quotumMap := make(map[time.Weekday]time.Duration)
	for weekday, dur := range input.Quota {
		duration, err := time.ParseDuration(dur)
		if err != nil {
			return model.Upkeep{}, err
		}
		quotumMap[time.Weekday(weekday)] = duration
	}

	var discounts []model.Discount
	for _, discountData := range input.Discounts {
		newDiscount := model.NewDiscount(discountData.Category, discountData.Argument)
		discounts = append(discounts, newDiscount)
	}

	upkeep := model.Upkeep{
		Version:    input.Version,
		Categories: infra.NewStackFromString(input.ActiveCategory),
		Quota:      quotumMap,
		Discounts:  discounts,
	}

	upkeep.Version = VERSION

	return upkeep, nil
}

func (r *UpkeepRepository) Insert(m model.Upkeep) error {
	quotumMap := make(map[int]string)
	for weekday, dur := range m.Quota {
		quotumMap[int(weekday)] = dur.String()
	}

	var discounts []discountJson
	for _, discount := range m.Discounts {
		discounts = append(discounts, discountJson{
			Category: discount.Category,
			Argument: discount.Argument,
		})
	}

	output := upkeepJson{
		Version:        m.Version,
		ActiveCategory: m.Categories.String(),
		Quota:          quotumMap,
		Discounts:      discounts,
	}

	if err := r.FileIO.Write("upkeep.json", output); err != nil {
		return err
	}

	return nil
}
