package core

import (
	"../utils"
	"gopkg.in/yaml.v2"
)

/*
	Program config struct
*/
type Config struct {
	Probes []Probe `yaml:"probes"`
}

var (
	// Program config
	DGConfig Config
)

/*
	Load program config from file
*/
func InitConfig(path string) {
	var content string // Config file content
	var err error      // Error handling

	// Read config file
	content, err = utils.FileReadAll(path)
	if err != nil {
		l.Critical("Content file read error:", err)
	}

	// Debug: display config file content
	l.Debug("Config file content:", "\n"+content)

	// Config file parsing: yaml => core.DGConfig
	err = yaml.Unmarshal([]byte(content), &DGConfig)
	if err != nil {
		l.Critical("error: %v", err)
	}

	l.Silly("DGConfig:\n", DGConfig)
}
