package infra

import (
	"fmt"
	"strings"
	"time"
)

type EscapeCode string

var (
	Reset  EscapeCode = "\033[0m"
	None   EscapeCode = ""
	Bold   EscapeCode = "\033[1m"
	Red    EscapeCode = "\033[31m"
	Green  EscapeCode = "\033[32m"
	Yellow EscapeCode = "\033[33m"
	White  EscapeCode = "\033[37m"

	BGGreen EscapeCode = "\033[42m"
)

type TerminalPrinter struct {
	lines []string
}

func (t *TerminalPrinter) PrintC(code EscapeCode, s string, args ...interface{}) *TerminalPrinter {
	t.addToBuffer(string(code)+s+string(Reset), args...)
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

func FormatDurationBracketed(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) - (hours * 60)

	return fmt.Sprintf("(%2d:%02d)", hours, minutes)
}

func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) - (hours * 60)

	return fmt.Sprintf("%d:%02d", hours, minutes)
}
