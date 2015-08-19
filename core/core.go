package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	dguard "github.com/90TechSAS/libgo-docker-guard"
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
	var resp *http.Response                     // Http response
	var req *http.Request                       // Http response
	var body []byte                             // Http body
	var err error                               // Error handling
	var containers map[string]*dguard.Container // Returned container list
	var probeID int                             // Probe ID

	probeID, err = GetProbeID(p.Name)
	if err != nil {
		l.Critical("MonitorProbe: Can't get probe ID", err)
	}

	// Reloading loop
	for {
		l.Verbose("Reloading", p.Name)

		// Make HTTP GET request
		reqURI := p.URI + "/list"
		l.Debug("GET", reqURI)
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
		err = json.Unmarshal([]byte(body), &containers)
		if err != nil {
			l.Error("MonitorProbe: Parsing container list:", err)
		}

		for _, c := range containers {
			var id int64
			var tmpContainer Container

			// Add containers in DB
			tmpContainer, err = GetContainerByCID(c.ID)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					sqlContainer := Container{0, c.ID, probeID, c.Hostname, c.Image, c.IPAddress, c.MacAddress}
					id, err = sqlContainer.Insert()
					if err != nil {
						l.Error("MonitorProbe: container insert:", err)
						continue
					}
				} else {
					l.Error("MonitorProbe: GetContainerById:", err)
					continue
				}
			} else {
				id = int64(tmpContainer.ID)
			}

			// Add stats in DB
			sqlStat := Stat{int(id), int64(c.Time), uint64(c.SizeRootFs), uint64(c.SizeRw), uint64(c.MemoryUsed), c.Running}
			err = sqlStat.Insert()
			if err != nil {
				l.Error("MonitorProbe: stat insert:", err)
				continue
			}
		}

		// time.Sleep(time.Second * time.Duration(p.ReloadTime))
	}
}
