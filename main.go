package main

import (
	"fmt"
	"os"
	"upkeep/app"
	"upkeep/app/view"
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
	mainRouter.Register("version", "show version", repository.Edit(app.HandleVersion))
	mainRouter.Register("start", "start a new block", repository.Edit(app.HandleStart))
	mainRouter.Register("stop", "stop the active block and save it", repository.Edit(app.HandleStop))
	mainRouter.Register("abort", "abort the active block without saving", repository.Edit(app.HandleAbort))
	mainRouter.Register("switch", "start a new block and put active category on the stack", repository.Edit(app.HandleSwitch))
	mainRouter.Register("continue", "start new block and pop active category from stack", repository.Edit(app.HandleContinue))
	mainRouter.Register("set", "set the active category", repository.Edit(app.HandleSet))
	mainRouter.Register("write", "write duration-only block", repository.Edit(app.HandleWrite))
	mainRouter.Register("edit/remove", "remove a block", repository.Edit(app.HandleRemove))
	mainRouter.Register("edit/update", "update block category", repository.Edit(app.HandleUpdate))
	mainRouter.Register("conf/quotum", "edit daily quotum", repository.Edit(app.HandleQuotum))
	mainRouter.Register("cat/quotum", "set the maximum daily quotum for a category", repository.Edit(app.HandleCategoryQuotum))
	mainRouter.Register("view/sheet", "view timesheet", repository.HandleView(view.ViewSheets))
	mainRouter.Register("view/day", "view totals by day", repository.HandleView(view.ViewDays))
	mainRouter.Register("view/cat", "view totals by category", repository.HandleView(view.ViewCategories))
	mainRouter.DefaultAction = "view/sheet"

	msg, err := mainRouter.Handle(infra.ParseArgs(os.Args))
	if err != nil {
		printer := infra.TerminalPrinter{}
		printer.Red("error: %s", err.Error()).Newline()
		fmt.Print(printer.String())
	}
	fmt.Printf("%s\n", msg)
}
