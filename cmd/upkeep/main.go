package main

import (
	"fmt"
	"os"

	"github.com/roelofruis/upkeep/internal/infra"
)

const (
	ModeProd  = "prod"
	ModeDev   = "dev"
	ModeDebug = "dbg"
)

// mode is set via ldflags in build
var mode = ModeProd

func main() {
	var homePath = "./data"
	prodMode := mode == ModeProd
	devMode := mode == ModeDev
	dbgMode := mode == ModeDebug
	if prodMode {
		var err error
		homePath, err = os.UserHomeDir()
		if err != nil {
			panic(err)
		}
	}

	var io infra.IO
	io = infra.FileIO{
		PrettyJson: devMode || dbgMode,
		HomePath:   homePath,
		DataFolder: ".upkeep",
	}

	if dbgMode {
		io = infra.IOLoggerDecorator{Inner: io}
	}

	router := Bootstrap(infra.SystemClock{}, io)

	msg, err := router.Handle(infra.ParseArgs(os.Args))
	if err != nil {
		printer := infra.TerminalPrinter{}
		printer.PrintC(infra.Red, "error: %s", err.Error()).Newline()
		fmt.Print(printer.String())
	}
	fmt.Printf("%s\n", msg)
}
