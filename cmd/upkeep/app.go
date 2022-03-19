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

	Router := infra.NewRouter()
	Router.Register("version", infra.Description{Base: "Show version information"}, app.Handle(HandleVersion))
	Router.Register("start", infra.Description{Base: "Start a new block", Extra: "Only first day of selection"}, app.Handle(HandleStart))
	Router.Register("stop", infra.Description{Base: "Stop the active block and save it", Extra: "Only first day of selection"}, app.Handle(HandleStop))
	Router.Register("abort", infra.Description{Base: "Abort the active block without saving", Extra: "Only first day of selection"}, app.Handle(HandleAbort))
	Router.Register("switch", infra.Description{Base: "Start a new block and put active category on the stack", Extra: "Only first day of selection"}, app.Handle(HandleSwitch))
	Router.Register("continue", infra.Description{Base: "Start new block and pop active category from stack", Extra: "Only first day of selection"}, app.Handle(HandleContinue))
	Router.Register("set", infra.Description{Base: "Set the active category"}, app.Handle(HandleSet))
	Router.Register("write", infra.Description{Base: "Write duration-only block", Extra: "Only first day of selection"}, app.Handle(HandleWrite))
	Router.Register("export", infra.Description{Base: "Write an export file to the current working directory"}, app.Handle(app.Export()))
	Router.Register("finalise", infra.Description{Base: "Mark timesheets as done, preventing further editing"}, app.Handle(HandleFinalise))
	Router.Register("unfinalise", infra.Description{Base: "Unmark timesheets as done"}, app.Handle(HandleUnfinalise))
	Router.Register("edit remove", infra.Description{Base: "Remove a block", Extra: "Only first day of selection"}, app.Handle(HandleRemove))
	Router.Register("edit restore", infra.Description{Base: "Restore a removed block", Extra: "Only first day of selection"}, app.Handle(HandleRestore))
	Router.Register("edit update", infra.Description{Base: "Update block category", Extra: "Only first day of selection"}, app.Handle(HandleUpdate))
	Router.Register("conf quotum", infra.Description{Base: "Edit daily quotum"}, app.Handle(HandleQuotum))
	Router.Register("cat quotum", infra.Description{Base: "Set the maximum daily quotum for a category"}, app.Handle(HandleCategoryQuotum))
	Router.Register("view sheet", infra.Description{Base: "View timesheet"}, app.Handle(ViewSheets))
	Router.Register("view day", infra.Description{Base: "View totals by day"}, app.Handle(ViewDays))
	Router.Register("view cat", infra.Description{Base: "View totals by category"}, app.Handle(ViewCategories))
	Router.DefaultAction = "view sheet"

	return Router
}
