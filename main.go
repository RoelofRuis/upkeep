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

	router.register("start", "start a new block", app.withDomain(handleStart))
	router.register("switch", "switch to a new block", app.withDomain(handleSwitch))
	router.register("stop", "stop the active block", app.withDomain(handleStop))
	router.register("tag", "change active tags", app.withDomain(handleTag))
	router.register("show", "show timesheet", app.withDomain(handleShow))
	router.register("purge", "purge timesheet", app.handlePurge)

	err, msg := router.handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
