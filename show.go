package main

import (
	"fmt"
	"strings"
	"timesheet/model"
)

func PrettyPrint(t *model.Timesheet) {
	if len(t.Blocks) == 0 {
		return
	}

	//dayStart := t.Blocks[0].Start
	//dayEnd := t.Blocks[len(t.Blocks)-1].End

	var availableHours = []int{9, 10, 11, 12, 13, 14, 15, 16}

	var hours []string
	var lines []string
	for _, hour := range availableHours {
		hours = append(hours, fmt.Sprintf("%02d:00", hour))
		lines = append(lines, fmt.Sprintf("|"))
	}

	fmt.Printf("%s\n", strings.Join(hours, " --- "))
	fmt.Printf("  %s\n", strings.Join(lines, "         "))
}
