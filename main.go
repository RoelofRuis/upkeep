package main

import (
	"fmt"
	"os"
	"timesheet/infra"
	"timesheet/repo"
)

type application struct {
	timekeepRepository  repo.TimekeepRepository
	timesheetRepository repo.TimesheetRepository
}

// mode is set via ldflags in build
var mode = "prod"

func main() {
	var homePath = "./dev-home"
	devMode := mode == "dev"
	if !devMode {
		var err error
		homePath, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
	}

	fileIO := infra.FileIO{PrettyJSON: devMode, HomePath: homePath}

	app := application{
		timekeepRepository:  repo.TimekeepRepository{FileIO: fileIO},
		timesheetRepository: repo.TimesheetRepository{FileIO: fileIO},
	}

	router := newRouter()
	router.register("test", "a test action", app.handleTest)

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
