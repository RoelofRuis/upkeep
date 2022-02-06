package infra

import (
	"fmt"
	"sort"
	"strings"
)

type Handler = func(args []string) (error, string)

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

func (r *Router) Handle(args []string) (error, string) {
	if len(args) == 0 {
		if r.DefaultAction == "" {
			return fmt.Errorf(r.NoCommandGivenMsg), r.HelpMessage()
		}

		args = append(args, r.DefaultAction)
	}

	if args[0] == "help" {
		return nil, r.HelpMessage()
	}

	h, has := r.actions[args[0]]
	if !has {
		return fmt.Errorf(r.NoMatchingHandlerMsg, args[0]), r.HelpMessage()
	}

	return h(args[1:])
}

func (r *Router) HelpMessage() string {
	lines := []string{fmt.Sprintf("available actions:")}

	actions := make([]string, 0, len(r.descriptions))
	for action := range r.descriptions {
		actions = append(actions, action)
	}
	sort.Strings(actions)

	for _, action := range actions {
		lines = append(lines, fmt.Sprintf("  - %s\n  %s", action, r.descriptions[action]))
	}
	return strings.Join(lines, "\n")
}
