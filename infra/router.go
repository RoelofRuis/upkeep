package infra

import (
	"fmt"
	"sort"
)

type Handler = func(params Params) (string, error)

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

func (r *Router) Handle(args Args) (string, error) {
	if args.Len() == 0 {
		if r.DefaultAction == "" {
			return r.HelpMessage(), fmt.Errorf(r.NoCommandGivenMsg)
		}

		args = args.Set([]string{r.DefaultAction})
	}

	if args.Path(1) == "help" {
		return r.HelpMessage(), nil
	}

	for i := args.Len(); i > 0; i-- {
		handler, has := r.actions[args.Path(i)]
		if has {
			return handler(args.GetParamsRemaining(i))
		}
	}

	return r.HelpMessage(), fmt.Errorf(r.NoMatchingHandlerMsg, args.Path(args.Len()))
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
