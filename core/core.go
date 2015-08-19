package core

import (
	"time"
)

var ()

type Probe struct {
	Name        string  `yaml:"name"`
	IPAddress   string  `yaml:"ip-address"`
	Port        string  `yaml:"port"`
	APIPassword string  `yaml:"api-password"`
	ReloadTime  float64 `yaml:"reload-time"`
}

/*
	Initialize Core
*/
func Init() {
	InitSQL()

	for _, probe := range DGConfig.Probes {
		go probe.MonitorProbe()
	}

	for {
		time.Sleep(time.Minute)
	}
}

/*
	Loop for monitoring a probe
*/
func (p *Probe) MonitorProbe() {
	for {
		l.Verbose("Reloading", p.Name, "("+p.IPAddress+":"+p.Port+")")
		// TODO

		time.Sleep(time.Second * time.Duration(p.ReloadTime))
	}
}
