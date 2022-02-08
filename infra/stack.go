package infra

import "strings"

type Stack []string

func NewStack() Stack {
	return Stack{}
}

func NewStackFromString(s string) Stack {
	if s == "" {
		return NewStack()
	}

	var elements Stack
	for _, elem := range strings.Split(s, "|") {
		elements = append(elements, elem)
	}

	return elements
}

func (ss Stack) IsEmpty() bool {
	return len(ss) == 0
}

func (ss Stack) Push(s string) Stack {
	return append(ss, s)
}

func (ss Stack) Pop() (Stack, string, bool) {
	if ss.IsEmpty() {
		return ss, "", false
	}

	index := len(ss) - 1
	element := ss[index]
	return ss[:index], element, true
}

func (ss Stack) Peek() string {
	if ss.IsEmpty() {
		return ""
	}

	index := len(ss) - 1
	element := ss[index]
	return element
}

func (ss Stack) String() string {
	return strings.Join(ss, "|")
}
