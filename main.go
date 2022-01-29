package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("help text here\n")
		return
	}

	tr := TimesheetRepository{path: "./data"}

	timesheet, err := tr.GetForDay(time.Now())
	if err != nil {
		panic(err) // TODO: cleaner error!
	}

	switch args[0] {
	case "start":
		err := timesheet.Start(time.Now())
		if err != nil {
			panic(err)
		}
	case "stop":
		err := timesheet.Stop(time.Now())
		if err != nil {
			panic(err)
		}
	case "show":
		PrettyPrint(timesheet)
	default:
		fmt.Printf("unknown command '%s'", args[0])
	}

	err = tr.Insert(timesheet)
	if err != nil {
		panic(err) // TODO: cleaner error!
	}

}
