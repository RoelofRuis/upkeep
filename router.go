package main

import "fmt"

type handler = func(args []string) (error, string)

type router struct {
	actions map[string]handler
}

func (r *router) register(action string, handler handler) {
	r.actions[action] = handler
}

func (r *router) handle(args []string) (error, string) {
	if len(args) == 0 {
		return fmt.Errorf("no command given"), ""
	}

	h, has := r.actions[args[0]]
	if !has {
		return fmt.Errorf("unknown command '%s'", args[0]), ""
	}

	return h(args[1:])
}
