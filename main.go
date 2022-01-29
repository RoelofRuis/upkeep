package main

import (
	"fmt"
	"time"
)

func main() {
	tr := TimesheetRepository{}

	timesheet, err := tr.Get()
	if err != nil {
		panic(err)
	}

	err = timesheet.Start(time.Now())
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", timesheet)

	err = tr.Insert(timesheet)
	if err != nil {
		panic(err)
	}

}
