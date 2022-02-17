package infra

import (
	"fmt"
	"strings"
	"time"
)

var (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	White  = "\033[37m"

	BGGreen = "\033[42m"
)

type TerminalPrinter struct {
	lines []string
}

func (t *TerminalPrinter) Red(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(Red+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) Green(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(Green+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) Yellow(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(Yellow+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) White(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(White+s+Reset, args...)
	return t
}

func (t *TerminalPrinter) BGGreen(s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(BGGreen+s+Reset, args...)
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

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) - (hours * 60)

	return fmt.Sprintf("(%2d:%02d)", hours, minutes)
}
