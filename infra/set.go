package infra

import "strings"

type Set []string

func NewSet() Set {
	return Set{}
}

func NewSetFromString(s string) Set {
	if s == "" {
		return NewSet()
	}

	var set Set
	for _, elem := range strings.Split(s, ",") {
		set = set.Add(elem)
	}

	return set
}

func (ss Set) IsEmpty() bool {
	return len(ss) == 0
}

func (ss Set) Add(elem string) Set {
	if ss.Contains(elem) {
		return ss
	}
	return append(ss, elem)
}

func (ss Set) Remove(elem string) Set {
	for i, e := range ss {
		if e == elem {
			ss[i] = ss[len(ss)-1]
			return ss[:len(ss)-1]
		}
	}
	return ss
}

func (ss Set) Contains(elem string) bool {
	for _, e := range ss {
		if e == elem {
			return true
		}
	}
	return false
}

func (ss Set) String() string {
	return strings.Join(ss, ",")
}
