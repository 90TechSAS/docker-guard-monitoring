package core

import (
	"os/exec"

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
	var err error  // Error handling
	var out []byte // Command output

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
