package core

import (
	"os"

	"github.com/nurza/logo"
)

var (
	// Logging
	l *logo.Logger
)

/*
	Exit the program when critical error
*/
func CriticalExit(s string) {
	println("Critical error => Exit program")
	os.Exit(1)
}

/*
	Initialize logger
*/
func InitLogger(l1, l2, l3 bool) *logo.Logger {
	var logger logo.Logger // Create logger
	var t *logo.Transport  // logger transport

	t = logger.AddTransport(logo.Console)              // Add a transport: Console
	t.AddColor(logo.ConsoleColor)                      // Add a color: Console color
	logger.AttachFunction(logo.Critical, CriticalExit) // Attach the function CriticalExit(string)
	logger.AddTime("[2006-01-02 15:04:05]")            // Add time prefix

	// Set verbose mode
	if l3 {
		logger.EnableLevel(logo.Silly)
		logger.EnableLevel(logo.Debug)
		logger.EnableLevel(logo.Verbose)
	} else if l2 {
		logger.EnableLevel(logo.Debug)
		logger.EnableLevel(logo.Verbose)
	} else if l1 {
		logger.EnableLevel(logo.Verbose)
	}
	logger.EnableLevel(logo.Info)
	logger.EnableLevel(logo.Warn)
	logger.EnableLevel(logo.Error)
	logger.EnableLevel(logo.Critical)

	l = &logger // Set core logger

	return l // return logger
}
