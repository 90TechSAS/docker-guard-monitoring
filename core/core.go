package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	HTTPClient *http.Client = &http.Client{}
)

type Probe struct {
	Name        string  `yaml:"name"`
	URI         string  `yaml:"uri"`
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
	var resp *http.Response // Http response
	var req *http.Request   // Http response
	var body []byte         // Http body
	var err error           // Error handling

	// Reloading loop
	for {
		l.Verbose("Reloading", p.Name)

		// Make HTTP GET request
		reqURI := p.URI + "/list"
		l.Verbose("GET", reqURI)
		req, err = http.NewRequest("GET", reqURI, bytes.NewBufferString(""))
		if err != nil {
			l.Error("MonitorProbe: Can't create", p.Name, "HTTP request:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}
		req.Header.Set("Auth", p.APIPassword)

		// Do request
		resp, err = HTTPClient.Do(req)
		if err != nil {
			l.Error("MonitorProbe: Can't get", p.Name, "container list:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}

		// Get request body
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			l.Error("MonitorProbe: Can't get", p.Name, "container list body:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}

		l.Silly("MonitorProbe:", "GET", reqURI, "body:\n", string(body))

		// Parse body
		// TODO

		time.Sleep(time.Second * time.Duration(p.ReloadTime))
	}
}
