package main

import (
	"fmt"
	"os"
)

type application struct {
	timesheetRepository TimesheetRepository
}

func main() {
	app := application{
		timesheetRepository: TimesheetRepository{path: "./data"},
	}

	router := router{actions: make(map[string]handler)}
	router.register("start", app.handleStart)
	router.register("stop", app.handleStop)
	router.register("tag", app.handleTag)
	router.register("show", app.handleShow)

	err, msg := router.handle(os.Args[1:])
	if err != nil {
		fmt.Printf("error: %s", err.Error())
	}
	fmt.Printf("%s\n", msg)
}
