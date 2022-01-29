package main

import (
	"os"
)

type application struct {
	timesheetRepository TimesheetRepository
}

func main() {
	app := application{
		timesheetRepository: TimesheetRepository{path: "./data"},
	}

	app.handle(os.Args[1:])
}
