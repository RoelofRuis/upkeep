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
	router.Register("version", infra.Description{Base: "Show version information"}, repository.Handle(app.HandleVersion))
	router.Register("start", infra.Description{Base: "Start a new block", Extra: "Only first day of selection"}, repository.Handle(app.HandleStart))
	router.Register("stop", infra.Description{Base: "Stop the active block and save it", Extra: "Only first day of selection"}, repository.Handle(app.HandleStop))
	router.Register("abort", infra.Description{Base: "Abort the active block without saving", Extra: "Only first day of selection"}, repository.Handle(app.HandleAbort))
	router.Register("switch", infra.Description{Base: "Start a new block and put active category on the stack", Extra: "Only first day of selection"}, repository.Handle(app.HandleSwitch))
	router.Register("continue", infra.Description{Base: "Start new block and pop active category from stack", Extra: "Only first day of selection"}, repository.Handle(app.HandleContinue))
	router.Register("set", infra.Description{Base: "Set the active category"}, repository.Handle(app.HandleSet))
	router.Register("write", infra.Description{Base: "Write duration-only block", Extra: "Only first day of selection"}, repository.Handle(app.HandleWrite))
	router.Register("export", infra.Description{Base: "Write an export file to the current working directory"}, repository.Handle(app.Export(fileIO)))
	router.Register("edit remove", infra.Description{Base: "Remove a block", Extra: "Only first day of selection"}, repository.Handle(app.HandleRemove))
	router.Register("edit restore", infra.Description{Base: "Restore a removed block", Extra: "Only first day of selection"}, repository.Handle(app.HandleRestore))
	router.Register("edit update", infra.Description{Base: "Update block category", Extra: "Only first day of selection"}, repository.Handle(app.HandleUpdate))
	router.Register("conf quotum", infra.Description{Base: "Edit daily quotum"}, repository.Handle(app.HandleQuotum))
	router.Register("cat quotum", infra.Description{Base: "Set the maximum daily quotum for a category"}, repository.Handle(app.HandleCategoryQuotum))
	router.Register("view sheet", infra.Description{Base: "View timesheet"}, repository.Handle(app.ViewSheets))
	router.Register("view day", infra.Description{Base: "View totals by day"}, repository.Handle(app.ViewDays))
	router.Register("view cat", infra.Description{Base: "View totals by category"}, repository.Handle(app.ViewCategories))
	router.DefaultAction = "view sheet"

	msg, err := router.Handle(infra.ParseArgs(os.Args))
	if err != nil {
		printer := infra.TerminalPrinter{}
		printer.PrintC(infra.Red, "error: %s", err.Error()).Newline()
		fmt.Print(printer.String())
	}
	fmt.Printf("%s\n", msg)
}
