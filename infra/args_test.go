package infra

import (
	"testing"
)

func TestNewParams(t *testing.T) {
	p := ParseArgs([]string{"upkeep", "view", "cat", "-dweek"})

	if p.ProgName != "upkeep" {
		t.Error("incorrect program name")
	}

	if p.Params.NamedParams["date"] != "week" {
		t.Error("incorrect duration range")
	}

	if p.Len() != 2 {
		t.Errorf("incorrect number of arguments")
	}

	if p.Path(2) != "view/cat" {
		t.Errorf("incorrect path exptected 'view/cat'")
	}

	if p.Path(1) != "view" {
		t.Errorf("incorrect path, expected 'view'")
	}

	if len(p.GetParamsRemaining(2).Params) != 0 {
		t.Errorf("incorrect remaining params")
	}

	if p.GetParamsRemaining(1).Params[0] != "cat" {
		t.Errorf("incorrect param")
	}
}
