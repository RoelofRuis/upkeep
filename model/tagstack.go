package model

import "strings"

type TagStack []TagSet

func NewTagStack() TagStack {
	return TagStack{}
}

func NewTagStackFromString(s string) TagStack {
	if s == "" {
		return NewTagStack()
	}

	var sets TagStack
	for _, tags := range strings.Split(s, "|") {
		sets = append(sets, NewTagSetFromString(tags))
	}

	return sets
}

func (ts TagStack) IsEmpty() bool {
	return len(ts) == 0
}

func (ts TagStack) Push(set TagSet) TagStack {
	return append(ts, set)
}

func (ts TagStack) Pop() (TagStack, TagSet, bool) {
	if ts.IsEmpty() {
		return ts, TagSet{}, false
	}

	index := len(ts) - 1
	element := ts[index]
	return ts[:index], element, true
}

func (ts TagStack) Peek() (TagSet, bool) {
	if ts.IsEmpty() {
		return TagSet{}, false
	}

	index := len(ts) - 1
	element := ts[index]
	return element, true
}

func (ts TagStack) String() string {
	var setStrings []string
	for _, set := range ts {
		setStrings = append(setStrings, set.String())
	}
	return strings.Join(setStrings, "|")
}
