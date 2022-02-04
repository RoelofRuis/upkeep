package infra

import (
	"fmt"
	"strings"
)

type Handler = func(args []string) (error, string)

type Router struct {
	actions      map[string]Handler
	descriptions map[string]string
}

func NewRouter() *Router {
	return &Router{
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
		return fmt.Errorf("no command given"), r.HelpMessage()
	}

	if args[0] == "help" {
		return nil, r.HelpMessage()
	}

	h, has := r.actions[args[0]]
	if !has {
		return fmt.Errorf("unknown command '%s'", args[0]), r.HelpMessage()
	}

	return h(args[1:])
}

func (r *Router) HelpMessage() string {
	lines := []string{fmt.Sprintf("available actions:")}
	for action, description := range r.descriptions {
		lines = append(lines, fmt.Sprintf("  - %s\n  %s", action, description))
	}
	return strings.Join(lines, "\n")
}
