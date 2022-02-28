package model

import (
	"fmt"
	"strings"
)

type Category struct {
	Group string
	Name  string
}

func NewCategoryFromString(s string) Category {
	parts := strings.Split(s, ".")
	if len(parts) < 2 {
		return Category{
			Group: "",
			Name:  parts[0],
		}
	}

	return Category{
		Group: parts[0],
		Name:  parts[1],
	}
}

func (c Category) GetName(byGroup bool) string {
	if !byGroup || c.Group == "" {
		return c.String()
	}

	return c.Group
}

func (c Category) String() string {
	if c.Group == "" {
		return c.Name
	}

	return fmt.Sprintf("%s.%s", c.Group, c.Name)
}
