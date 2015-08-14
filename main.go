package main

import (
	"flag"
	"fmt"
	"os"

	"./core"

	"github.com/nurza/logo"
)

var (
	// Logging
	l *logo.Logger

	DisplayHelp1 = flag.Bool("h", false, "Display help")
	DisplayHelp2 = flag.Bool("help", false, "Display help")
	ConfigFile   = flag.String("f", "config.yaml", "Configuration file")
	Verbose1     = flag.Bool("v", false, "Verbose mode 1 : display verbose logs")
	Verbose2     = flag.Bool("vv", false, "Verbose mode 2 : display verbose and debug logs")
	Verbose3     = flag.Bool("vvv", false, "Verbose mode 3 : display verbose, debug and silly logs")
)

/*
	Help message
*/
var help string = `Usage of Docker Guard System Monitoring: dgs-monitoring <options>
	Options (<option>=[default value]):
		-f="config.yaml": Configuration file
		-h / --help: Display help
		-v: Verbose mode 1 : display verbose logs
		-vv: Verbose mode 2 : display verbose and debug logs
		-vvv: Verbose mode 3 : display verbose, debug and silly logs`

func main() {
	// Flags
	flag.Parse()
	// If option help, display help message and exit
	if *DisplayHelp1 || *DisplayHelp2 {
		fmt.Println(help)
		os.Exit(0)
	}

	// Logging
	println("Init logger...")
	l = core.InitLogger(*Verbose1, *Verbose2, *Verbose3)
	l.Verbose("logger OK")

	// Config
	l.Verbose("Init config")
	core.InitConfig(*ConfigFile)
	l.Verbose("config OK")

	// Init Core
	core.Init()

}
