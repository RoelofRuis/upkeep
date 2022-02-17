package infra

import (
	"regexp"
	"strconv"
	"strings"
)

type Args struct {
	ProgName string
	Args     []string
	Params   Params
}

var namedArgument = regexp.MustCompile("^([a-z]+):([a-zA-Z0-9_-]+)$")

func ParseArgs(a []string) Args {
	args := Args{
		ProgName: a[0],
		Params:   Params{NamedParams: map[string]string{"d": "day"}},
	}

	for _, s := range a[1:] {
		matches := namedArgument.FindStringSubmatch(s)
		if len(matches) == 3 {
			args.Params.NamedParams[matches[1]] = matches[2]
			continue
		}

		args.Args = append(args.Args, s)
	}

	return args
}

func (a Args) Set(args []string) Args {
	a.Args = args
	return a
}

func (a Args) GetParamsRemaining(len int) Params {
	params := a.Params
	params.Params = a.Args[len:]
	return params
}

func (a Args) Len() int {
	return len(a.Args)
}

func (a Args) Path(len int) string {
	return strings.Join(a.Args[:len], " ")
}

type Params struct {
	Params      []string
	NamedParams map[string]string
}

func (p Params) Len() int {
	return len(p.Params)
}

func (p Params) Get(index int) string {
	if index+1 > len(p.Params) {
		return ""
	}
	return p.Params[index]
}

func (p Params) GetInt(index int) (int, error) {
	i, err := strconv.ParseInt(p.Get(index), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i), nil
}

func (p Params) GetNamed(name string, fallback string) string {
	val, has := p.NamedParams[name]
	if !has {
		return fallback
	}
	return val
}
