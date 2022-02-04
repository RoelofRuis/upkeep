package main

import (
	"fmt"
	"os"
	"timesheet/app"
	"timesheet/infra"
	"timesheet/model/repo"
)

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
		DebugEnabled: devMode,
		HomePath:     homePath,
		DataFolder:   ".upkeep",
	}

	repository := app.Repository(repo.New(fileIO))

	router := infra.NewRouter("show")
	router.Register("start", "start a new block", repository.Edit(app.HandleStart))
	router.Register("switch", "switch to a new block", repository.Edit(app.HandleSwitch))
	router.Register("stop", "stop the active block", repository.Edit(app.HandleStop))
	router.Register("tag", "change active tags", repository.Edit(app.HandleTag))
	router.Register("purge", "purge timesheet", repository.Edit(app.HandlePurge))

	router.Register("show", "show timesheet", repository.Edit(app.HandleShow))

	err, msg := router.Handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
