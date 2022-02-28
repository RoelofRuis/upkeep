package model

type CategorySetting struct {
	Name         string
	MaxDayQuotum Duration
}

func NewCategorySetting(name string) CategorySetting {
	return CategorySetting{
		Name:         name,
		MaxDayQuotum: NewDuration(),
	}
}

type CategorySettings []CategorySetting

func (c CategorySettings) IsEmpty() bool {
	return len(c) == 0
}

func (c CategorySettings) Add(cat CategorySetting) CategorySettings {
	for i, e := range c {
		if e.Name == cat.Name {
			c[i] = cat
			return c
		}
	}

	return append(c, cat)
}

func (c CategorySettings) Remove(elem string) CategorySettings {
	for i, e := range c {
		if e.Name == elem {
			c[i] = c[len(c)-1]
			return c[:len(c)-1]
		}
	}
	return c
}

func (c CategorySettings) Contains(elem string) bool {
	for _, e := range c {
		if e.Name == elem {
			return true
		}
	}
	return false
}

func (c CategorySettings) Get(name string) CategorySetting {
	for _, cat := range c {
		if cat.Name == name {
			return cat
		}
	}
	return NewCategorySetting(name)
}

func (c CategorySettings) Names() []string {
	var list []string
	for _, cat := range c {
		list = append(list, cat.Name)
	}
	return list
}
