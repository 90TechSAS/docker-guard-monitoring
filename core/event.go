package core

import (
	"os/exec"
	"regexp"

	"../utils"

	dguard "github.com/90TechSAS/libgo-docker-guard"
)

var ()

type Transport struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

func EventController() {
	// TODO
}

func Alert(event dguard.Event) {
	var err error          // Error handling
	var out []byte         // Command output
	var alert bool = false // true if the alert needs to be sent

	for _, rgxp := range DGConfig.DockerGuard.Event.Watch {
		r, err := regexp.MatchString(rgxp, event.Target)
		if err != nil {
			l.Error("Error processing regexp:", err)
			return
		}
		if r {
			alert = true
			break
		}
	}

	if !alert {
		return
	}

	// Exec transports
	for _, t := range DGConfig.DockerGuard.Event.Transports {
		out, err = exec.Command(t.Path,
			utils.I2S(event.Severity),
			event.TypeToString(),
			event.Target,
			event.Probe,
			event.Data).Output()
		if err != nil {
			l.Error("Error transport ("+t.Name+") Out:", string(out))
			return
		}
		l.Debug("Transport ("+t.Name+") Out:", string(out))
	}
}
