package main

import (
	"fmt"
	"strings"
)

type handler = func(args []string) (error, string)

type router struct {
	actions      map[string]handler
	descriptions map[string]string
}

func newRouter() *router {
	return &router{
		actions:      make(map[string]handler),
		descriptions: make(map[string]string),
	}
}

func (r *router) register(action string, description string, handler handler) {
	r.actions[action] = handler
	r.descriptions[action] = description
}

func (r *router) handle(args []string) (error, string) {
	if len(args) == 0 {
		return fmt.Errorf("no command given"), r.helpMessage()
	}

	h, has := r.actions[args[0]]
	if !has {
		return fmt.Errorf("unknown command '%s'", args[0]), r.helpMessage()
	}

	return h(args[1:])
}

func (r *router) helpMessage() string {
	lines := []string{fmt.Sprintf("available actions:")}
	for action, description := range r.descriptions {
		lines = append(lines, fmt.Sprintf("  - %s\n  %s", action, description))
	}
	return strings.Join(lines, "\n")
}
