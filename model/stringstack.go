package model

import "strings"

type StringStack []string

func NewStringStack() StringStack {
	return StringStack{}
}

func NewStringStackFromString(s string) StringStack {
	if s == "" {
		return NewStringStack()
	}

	var elements StringStack
	for _, elem := range strings.Split(s, "|") {
		elements = append(elements, elem)
	}

	return elements
}

func (ss StringStack) IsEmpty() bool {
	return len(ss) == 0
}

func (ss StringStack) Push(s string) StringStack {
	return append(ss, s)
}

func (ss StringStack) Pop() (StringStack, string, bool) {
	if ss.IsEmpty() {
		return ss, "", false
	}

	index := len(ss) - 1
	element := ss[index]
	return ss[:index], element, true
}

func (ss StringStack) Peek() string {
	if ss.IsEmpty() {
		return ""
	}

	index := len(ss) - 1
	element := ss[index]
	return element
}

func (ss StringStack) String() string {
	return strings.Join(ss, "|")
}
