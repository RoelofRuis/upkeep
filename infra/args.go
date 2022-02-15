package infra

import (
	"fmt"
	"strconv"
	"strings"
)

type Args struct {
	ProgName string
	Args     []string
	Params   Params
}

func ParseArgs(a []string) Args {
	args := Args{
		ProgName: a[0],
		Params:   Params{NamedParams: map[string]string{"date": "today"}},
	}

	for _, s := range a[1:] {
		if strings.HasPrefix(s, "-d") {
			args.Params.NamedParams["date"] = strings.TrimPrefix(s, "-d")
		} else {
			args.Args = append(args.Args, s)
		}
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
	return strings.Join(a.Args[:len], "/")
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

func (p Params) GetNamed(name string) (string, error) {
	val, has := p.NamedParams[name]
	if !has {
		return "", fmt.Errorf("no named parameter '%s'", name)
	}
	return val, nil
}
