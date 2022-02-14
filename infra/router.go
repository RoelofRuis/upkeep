package infra

import (
	"fmt"
	"sort"
)

type Handler = func(args []string) (string, error)

type Router struct {
	NoCommandGivenMsg    string
	NoMatchingHandlerMsg string
	DefaultAction        string
	HelpActive           bool

	actions      map[string]Handler
	descriptions map[string]string
}

func NewRouter() *Router {
	return &Router{
		DefaultAction:        "",
		NoCommandGivenMsg:    "no command given",
		NoMatchingHandlerMsg: "unknown command '%s'",

		actions:      make(map[string]Handler),
		descriptions: make(map[string]string),
	}
}

func (r *Router) Register(action string, description string, handler Handler) {
	r.actions[action] = handler
	r.descriptions[action] = description
}

func (r *Router) Handle(args []string) (string, error) {
	if len(args) == 0 {
		if r.DefaultAction == "" {
			return r.HelpMessage(), fmt.Errorf(r.NoCommandGivenMsg)
		}

		args = append(args, r.DefaultAction)
	}

	if args[0] == "help" {
		return r.HelpMessage(), nil
	}

	h, has := r.actions[args[0]]
	if !has {
		return r.HelpMessage(), fmt.Errorf(r.NoMatchingHandlerMsg, args[0])
	}

	return h(args[1:])
}

func (r *Router) HelpMessage() string {
	printer := TerminalPrinter{}
	printer.Print("available actions:").Newline()

	actions := make([]string, 0, len(r.descriptions))
	for action := range r.descriptions {
		actions = append(actions, action)
	}
	sort.Strings(actions)

	for _, action := range actions {
		printer.Print("> ").
			Bold("%s", action).
			Newline().
			White("  %s", r.descriptions[action]).
			Newline()
	}
	return printer.String()
}
