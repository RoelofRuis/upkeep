package infra

import (
	"fmt"
	"strings"
)

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Yellow = "\033[33m"
)

type TerminalPrinter struct {
	lines []string
}

func (t *TerminalPrinter) Red(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(Red+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) Yellow(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(Yellow+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) Bold(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(Bold+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) Print(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(s, args...)
	return t
}

func (t *TerminalPrinter) Newline() *TerminalPrinter {
	t.lines = append(t.lines, "")
	return t
}

func (t *TerminalPrinter) addToBuffer(s string, args ...interface{}) {
	text := fmt.Sprintf(s, args...)

	if len(t.lines) == 0 {
		t.lines = append(t.lines, "")
	}
	t.lines[len(t.lines)-1] += text
}

func (t *TerminalPrinter) String() string {
	return strings.Join(t.lines, "\n")
}
