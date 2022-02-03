package main

import (
	"fmt"
	"os"
)

type application struct {
	fileIO              FileIO
	timesheetRepository TimesheetRepository
}

// mode is set via ldflags in build
var mode = "prod"

func main() {
	var path = "./dev-home"
	if mode == "prod" {
		var err error
		path, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
	}

	fileIO := FileIO{path: path}

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
	router.register("break", "start a break block", app.handleBreak)

	err, msg := router.handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
