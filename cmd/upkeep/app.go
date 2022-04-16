package main

import (
	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model/repo"
)

type App struct {
	Clock infra.Clock
	IO    infra.IO
	Repo  repo.Repository
}

func Bootstrap(
	clock infra.Clock,
	io infra.IO,
) *infra.Router {
	app := App{
		Clock: clock,
		IO:    io,
		Repo:  repo.New(io),
	}

	router := infra.NewRouter()
	router.Register("version", infra.Description{Base: "Show version information"}, app.Handle(HandleVersion))
	router.Register("start", infra.Description{Base: "Start a new block", Extra: "Only first day of selection"}, app.Handle(HandleStart))
	router.Register("stop", infra.Description{Base: "Stop the active block and save it", Extra: "Only first day of selection"}, app.Handle(HandleStop))
	router.Register("abort", infra.Description{Base: "Abort the active block without saving", Extra: "Only first day of selection"}, app.Handle(HandleAbort))
	router.Register("switch", infra.Description{Base: "Start a new block and put active category on the stack", Extra: "Only first day of selection"}, app.Handle(HandleSwitch))
	router.Register("swap", infra.Description{Base: "Swap last and second last categories on the stack", Extra: "Only first day of selection"}, app.Handle(HandleSwap))
	router.Register("continue", infra.Description{Base: "Start new block and pop active category from stack", Extra: "Only first day of selection"}, app.Handle(HandleContinue))
	router.Register("set", infra.Description{Base: "Set the active category"}, app.Handle(HandleSet))
	router.Register("write", infra.Description{Base: "Write duration-only block", Extra: "Only first day of selection"}, app.Handle(HandleWrite))
	router.Register("export", infra.Description{Base: "Write an export file to the current working directory"}, app.Handle(app.Export()))
	router.Register("finalise", infra.Description{Base: "Mark timesheets as done, preventing further editing"}, app.Handle(HandleFinalise))
	router.Register("unfinalise", infra.Description{Base: "Unmark timesheets as done"}, app.Handle(HandleUnfinalise))
	router.Register("edit remove", infra.Description{Base: "Remove a block", Extra: "Only first day of selection"}, app.Handle(HandleRemove))
	router.Register("edit restore", infra.Description{Base: "Restore a removed block", Extra: "Only first day of selection"}, app.Handle(HandleRestore))
	router.Register("edit update", infra.Description{Base: "Update block category", Extra: "Only first day of selection"}, app.Handle(HandleUpdate))
	router.Register("conf quotum", infra.Description{Base: "Edit daily quotum"}, app.Handle(HandleQuotum))
	router.Register("cat quotum", infra.Description{Base: "Set the maximum daily quotum for a category"}, app.Handle(HandleCategoryQuotum))
	router.Register("view sheet", infra.Description{Base: "View timesheet"}, app.Handle(ViewSheets))
	router.Register("view day", infra.Description{Base: "View totals by day"}, app.Handle(ViewDays))
	router.Register("view cat", infra.Description{Base: "View totals by category"}, app.Handle(ViewCategories))
	router.DefaultAction = "view sheet"

	return router
}
