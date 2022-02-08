package model

import "strings"

type TagSet struct {
	tags []string
}

func NewTagSet() TagSet {
	return TagSet{tags: nil}
}

func NewTagSetFromString(s string) TagSet {
	if s == "" {
		return NewTagSet()
	}

	return TagSet{tags: strings.Split(s, ",")}
}

func (ts TagSet) Add(tag string) TagSet {
	if ts.Contains(tag) {
		return ts
	}
	return TagSet{tags: append(ts.tags, tag)}
}

func (ts TagSet) Remove(tag string) TagSet {
	for i, t := range ts.tags {
		if t == tag {
			ts.tags[i] = ts.tags[len(ts.tags)-1]
			ts.tags = ts.tags[:len(ts.tags)-1]
			return TagSet{tags: ts.tags}
		}
	}
	return ts
}

func (ts TagSet) Contains(tag string) bool {
	for _, t := range ts.tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (ts TagSet) String() string {
	return strings.Join(ts.tags, ",")
}
