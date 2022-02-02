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

	router := router{actions: make(map[string]handler)}
	router.register("start", app.handleStart)
	router.register("stop", app.handleStop)
	router.register("tag", app.handleTag)
	router.register("show", app.handleShow)
	router.register("purge", app.handlePurge)

	err, msg := router.handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
