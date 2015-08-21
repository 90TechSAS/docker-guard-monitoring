package core

import (
	dguard "github.com/90TechSAS/libgo-docker-guard"
)

var ()

type Event dguard.Event

type Transport struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Ip   string `yaml:"ip"`
	Port string `yaml:"port"`
}

func EventController() {
	// TODO
}

func (e *Event) Alert() {
	// TODO
}
