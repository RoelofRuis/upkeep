package model

type Categories []Category

func (c Categories) IsEmpty() bool {
	return len(c) == 0
}

func (c Categories) Add(cat Category) Categories {
	for i, e := range c {
		if e.Name == cat.Name {
			c[i] = cat
			return c
		}
	}

	return append(c, cat)
}

func (c Categories) Remove(elem string) Categories {
	for i, e := range c {
		if e.Name == elem {
			c[i] = c[len(c)-1]
			return c[:len(c)-1]
		}
	}
	return c
}

func (c Categories) Contains(elem string) bool {
	for _, e := range c {
		if e.Name == elem {
			return true
		}
	}
	return false
}

func (c Categories) Get(name string) Category {
	for _, cat := range c {
		if cat.Name == name {
			return cat
		}
	}
	return NewCategory(name)
}

func (c Categories) Names() []string {
	var list []string
	for _, cat := range c {
		list = append(list, cat.Name)
	}
	return list
}
