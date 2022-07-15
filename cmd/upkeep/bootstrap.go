package main

import (
	"github.com/roelofruis/upkeep/internal/infra"
	"github.com/roelofruis/upkeep/internal/model/repo"
)

type Dependencies struct {
	Clock infra.Clock
	IO    infra.IO
	Repo  repo.Repository
}

func Bootstrap(
	clock infra.Clock,
	io infra.IO,
) *infra.Router {
	dependencies := &Dependencies{
		Clock: clock,
		IO:    io,
		Repo:  repo.New(io),
	}

	router := infra.NewRouter()
	router.Register("version", infra.Description{Base: "Show version information"}, Handle(dependencies, HandleVersion))
	router.Register("start", infra.Description{Base: "Start a new block", Extra: "Only first day of selection"}, Handle(dependencies, HandleStart))
	router.Register("stop", infra.Description{Base: "Stop the active block and save it", Extra: "Only first day of selection"}, Handle(dependencies, HandleStop))
	router.Register("abort", infra.Description{Base: "Abort the active block without saving", Extra: "Only first day of selection"}, Handle(dependencies, HandleAbort))
	router.Register("switch", infra.Description{Base: "Start a new block and put active category on the stack", Extra: "Only first day of selection"}, Handle(dependencies, HandleSwitch))
	router.Register("swap", infra.Description{Base: "Swap last and second last categories on the stack", Extra: "Only first day of selection"}, Handle(dependencies, HandleSwap))
	router.Register("continue", infra.Description{Base: "Start new block and pop active category from stack", Extra: "Only first day of selection"}, Handle(dependencies, HandleContinue))
	router.Register("set", infra.Description{Base: "Set the active category"}, Handle(dependencies, HandleSet))
	router.Register("write", infra.Description{Base: "Write duration-only block", Extra: "Only first day of selection"}, Handle(dependencies, HandleWrite))
	router.Register("export", infra.Description{Base: "Write an export file to the current working directory"}, Handle(dependencies, dependencies.Export()))
	router.Register("finalise", infra.Description{Base: "Mark timesheets as done, preventing further editing"}, Handle(dependencies, HandleFinalise))
	router.Register("unfinalise", infra.Description{Base: "Unmark timesheets as done"}, Handle(dependencies, HandleUnfinalise))
	router.Register("quotum", infra.Description{Base: "Adjust quotum", Extra: "Only first day of selection"}, Handle(dependencies, HandleQuotum))
	router.Register("edit remove", infra.Description{Base: "Remove a block", Extra: "Only first day of selection"}, Handle(dependencies, HandleRemove))
	router.Register("edit restore", infra.Description{Base: "Restore a removed block", Extra: "Only first day of selection"}, Handle(dependencies, HandleRestore))
	router.Register("edit update", infra.Description{Base: "Update block category", Extra: "Only first day of selection"}, Handle(dependencies, HandleUpdate))
	router.Register("conf quotum", infra.Description{Base: "Edit daily quotum"}, Handle(dependencies, HandleDayQuotum))
	router.Register("cat quotum", infra.Description{Base: "Set the maximum daily quotum for a category"}, Handle(dependencies, HandleCategoryQuotum))
	router.Register("view sheet", infra.Description{Base: "View timesheet"}, Handle(dependencies, ViewSheets))
	router.Register("view day", infra.Description{Base: "View totals by day"}, Handle(dependencies, ViewDays))
	router.Register("view cat", infra.Description{Base: "View totals by category"}, Handle(dependencies, ViewCategories))
	router.DefaultAction = "view sheet"

	return router
}
