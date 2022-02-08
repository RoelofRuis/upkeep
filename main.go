package main

import (
	"fmt"
	"os"
	"upkeep/app"
	"upkeep/infra"
	"upkeep/model/repo"
)

const (
	ModeProd  = "prod"
	ModeDev   = "dev"
	ModeDebug = "dbg"
)

// mode is set via ldflags in build
var mode = ModeProd

func main() {
	var homePath = "./dev-home"
	prodMode := mode == ModeProd
	devMode := mode == ModeDev
	dbgMode := mode == ModeDebug
	if prodMode {
		var err error
		homePath, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
	}

	fileIO := infra.FileIO{
		PrettyJson:   devMode || dbgMode,
		DebugEnabled: dbgMode,
		HomePath:     homePath,
		DataFolder:   ".upkeep",
	}

	repository := app.Repository(repo.New(fileIO))

	mainRouter := infra.NewRouter()
	mainRouter.Register("start", "start a new block", repository.Edit(app.HandleStart))
	mainRouter.Register("stop", "stop the active block and save it", repository.Edit(app.HandleStop))
	mainRouter.Register("abort", "abort the active block without saving", repository.Edit(app.HandleAbort))
	mainRouter.Register("switch", "switch to a new block with new category", repository.Edit(app.HandleSwitch))
	mainRouter.Register("continue", "switch back to the old category", repository.Edit(app.HandleContinue))
	mainRouter.Register("cat", "change active category", repository.Edit(app.HandleCategory))

	mainRouter.Register("remove", "remove a time block", repository.Edit(app.HandleRemove))

	confRouter := infra.NewRouter()
	mainRouter.Register("conf", "edit configuration values", confRouter.Handle)
	confRouter.Register("quotum", "edit daily quotum", repository.Edit(app.HandleQuotum))

	mainRouter.Register("day", "show timesheet for today", repository.Edit(app.HandleDay))
	mainRouter.DefaultAction = "day"

	err, msg := mainRouter.Handle(os.Args[1:])
	if err != nil {
		printer := infra.TerminalPrinter{}
		printer.Red("error: %s", err.Error()).Newline()
		fmt.Print(printer.String())
	}
	fmt.Printf("%s\n", msg)
}
