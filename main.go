package main

import (
	"fmt"
	"os"
	"timesheet/infra"
	"timesheet/repo"
)

type application struct {
	upkeepRepository    repo.UpkeepRepository
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

	fileIO := infra.FileIO{
		PrettyJSON: devMode,
		HomePath:   homePath,
		DataFolder: ".upkeep",
	}

	app := application{
		upkeepRepository:    repo.UpkeepRepository{FileIO: fileIO},
		timesheetRepository: repo.TimesheetRepository{FileIO: fileIO},
	}

	router := newRouter()

	router.register("start", "start a new block", app.withModel(handleStart))
	router.register("stop", "stop the active block", app.withModel(handleStop))
	router.register("tag", "change active tags", app.withModel(handleTag))
	router.register("show", "show timesheet", app.withModel(handleShow))
	router.register("purge", "purge timesheet", app.handlePurge)

	err, msg := router.handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
