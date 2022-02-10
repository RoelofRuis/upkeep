package model

import "time"

type Discount struct {
	Category string
	Argument string
}

func NewDiscount(cat string, arg string) Discount {
	return Discount{
		Category: cat,
		Argument: arg,
	}
}

func DiscountAll(cat string) Discount {
	return NewDiscount(cat, "")
}

func DiscountMax(cat string, dur time.Duration) Discount {
	return NewDiscount(cat, dur.Truncate(time.Second).String())
}

func (d Discount) Measure() *Discounter {
	if d.Argument == "" {
		return &Discounter{
			cat:          d.Category,
			durRemaining: time.Duration(0),
		}
	}
	dur, _ := time.ParseDuration(d.Argument)
	return &Discounter{
		cat:          d.Category,
		durRemaining: dur,
	}

}

type Discounter struct {
	cat          string
	durRemaining time.Duration
}

func (m *Discounter) GetTimeRemaining(cat string, d time.Duration) time.Duration {
	if cat != m.cat {
		return d
	}

	if m.durRemaining > d {
		m.durRemaining -= d
		return d
	}

	remaining := m.durRemaining
	m.durRemaining = 0
	return remaining
}
