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

	router := infra.NewRouter()
	router.Register("version", "show version", repository.Handle(app.HandleVersion))
	router.Register("start", "start a new block", repository.Handle(app.HandleStart))
	router.Register("stop", "stop the active block and save it", repository.Handle(app.HandleStop))
	router.Register("abort", "abort the active block without saving", repository.Handle(app.HandleAbort))
	router.Register("switch", "start a new block and put active category on the stack", repository.Handle(app.HandleSwitch))
	router.Register("continue", "start new block and pop active category from stack", repository.Handle(app.HandleContinue))
	router.Register("set", "set the active category", repository.Handle(app.HandleSet))
	router.Register("write", "write duration-only block", repository.Handle(app.HandleWrite))
	router.Register("export", "write an export file", repository.Handle(app.Export))
	router.Register("edit remove", "remove a block", repository.Handle(app.HandleRemove))
	router.Register("edit restore", "restore a removed block", repository.Handle(app.HandleRestore))
	router.Register("edit update", "update block category", repository.Handle(app.HandleUpdate))
	router.Register("conf quotum", "edit daily quotum", repository.Handle(app.HandleQuotum))
	router.Register("cat quotum", "set the maximum daily quotum for a category", repository.Handle(app.HandleCategoryQuotum))
	router.Register("view sheet", "view timesheet", repository.Handle(app.ViewSheets))
	router.Register("view day", "view totals by day", repository.Handle(app.ViewDays))
	router.Register("view cat", "view totals by category", repository.Handle(app.ViewCategories))
	router.DefaultAction = "view sheet"

	msg, err := router.Handle(infra.ParseArgs(os.Args))
	if err != nil {
		printer := infra.TerminalPrinter{}
		printer.PrintC(infra.Red, "error: %s", err.Error()).Newline()
		fmt.Print(printer.String())
	}
	fmt.Printf("%s\n", msg)
}
