package main

import (
	"fmt"
	"os"
)

type application struct {
	fileIO              FileIO
	timesheetRepository TimesheetRepository
}

func main() {
	fileIO := FileIO{path: "./data"}

	app := application{
		fileIO:              fileIO,
		timesheetRepository: TimesheetRepository{fileIO: fileIO},
	}

	router := newRouter()
	router.register("start", "start a new block", app.handleStart)
	router.register("stop", "stop the active block", app.handleStop)
	router.register("tag", "change active tags", app.handleTag)
	router.register("show", "show timesheet", app.handleShow)
	router.register("purge", "purge timesheet", app.handlePurge)

	err, msg := router.handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
