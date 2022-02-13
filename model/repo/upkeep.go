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
	Version            string                 `json:"version"`
	SelectedCategories string                 `json:"selected_categories"`
	Quota              map[int]model.Duration `json:"quota"`
	Categories         []categoryJson         `json:"categories"`
}

type categoryJson struct {
	Name         string         `json:"name"`
	MaxDayQuotum model.Duration `json:"max_day_quotum"`
}

const VERSION = "0.2"

func (r *UpkeepRepository) Get() (model.Upkeep, error) {
	input := upkeepJson{
		Version: "0.1",
	}

	if err := r.FileIO.Read("upkeep.json", &input); err != nil {
		return model.Upkeep{}, err
	}

	quotumMap := make(map[time.Weekday]model.Duration)
	for weekday, dur := range input.Quota {
		quotumMap[time.Weekday(weekday)] = dur
	}

	var categories model.Categories
	for _, categoryData := range input.Categories {
		newCategory := model.NewCategory(categoryData.Name)
		newCategory.MaxDayQuotum = categoryData.MaxDayQuotum
		categories = append(categories, newCategory)
	}

	upkeep := model.Upkeep{
		Version:            input.Version,
		SelectedCategories: infra.NewStackFromString(input.SelectedCategories),
		Quota:              quotumMap,
		Categories:         categories,
	}

	upkeep.Version = VERSION

	return upkeep, nil
}

func (r *UpkeepRepository) Insert(m model.Upkeep) error {
	quotumMap := make(map[int]model.Duration)
	for weekday, dur := range m.Quota {
		quotumMap[int(weekday)] = dur
	}

	var categories []categoryJson
	for _, category := range m.Categories {
		categories = append(categories, categoryJson{
			Name:         category.Name,
			MaxDayQuotum: category.MaxDayQuotum,
		})
	}

	output := upkeepJson{
		Version:            m.Version,
		SelectedCategories: m.SelectedCategories.String(),
		Quota:              quotumMap,
		Categories:         categories,
	}

	if err := r.FileIO.Write("upkeep.json", output); err != nil {
		return err
	}

	return nil
}
