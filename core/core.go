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
	// HTTP client used to get probe infos
	HTTPClient *http.Client = &http.Client{}
)

/*
	Probe
*/
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
	// Init SQL client
	InitSQL()

	// Launch probe monitors
	for _, probe := range DGConfig.Probes {
		go MonitorProbe(probe)
	}

	// Launch event controller
	EventController()

	// Launch API
	HTTPServer()
}

/*
	Loop for monitoring a probe
*/
func MonitorProbe(p Probe) {
	var resp *http.Response                        // Http response
	var req *http.Request                          // Http response
	var body []byte                                // Http body
	var err error                                  // Error handling
	var containers map[string]*dguard.Container    // Returned container list
	var oldContainers map[string]*dguard.Container // Old returned container list (used to compare running state)
	var dbContainers []Container                   // Containers in DB
	var probeID int                                // Probe ID

	// Get probe ID in DB, create it if does not exists
	probeID, err = GetProbeID(p.Name)
	if err != nil {
		l.Critical("MonitorProbe ("+p.Name+"): Can't get probe ID", err)
	}

	// Reloading loop
	for {
		oldContainers = containers
		containers = nil
		l.Verbose("Reloading", p.Name)

		// Make HTTP GET request
		reqURI := p.URI + "/list"
		l.Debug("GET", reqURI)
		req, err = http.NewRequest("GET", reqURI, bytes.NewBufferString(""))
		if err != nil {
			l.Error("MonitorProbe ("+p.Name+"): Can't create", p.Name, "HTTP request:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}
		req.Header.Set("Auth", p.APIPassword)

		// Do request
		resp, err = HTTPClient.Do(req)
		if err != nil {
			l.Error("MonitorProbe ("+p.Name+"): Can't get", p.Name, "container list:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}

		// Get request body
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			l.Error("MonitorProbe ("+p.Name+"): Can't get", p.Name, "container list body:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}

		l.Silly("MonitorProbe ("+p.Name+"):", "GET", reqURI, "body:\n", string(body))

		// Parse body
		err = json.Unmarshal([]byte(body), &containers)
		if err != nil {
			l.Error("MonitorProbe ("+p.Name+"): Parsing container list:", err)
		}

		// Remove in DB old removed containers
		dbContainers, err = GetContainersBy("probeid", probeID)
		if err != nil {
			l.Error("MonitorProbe ("+p.Name+"): Can't get", p.Name, "container list in DB:", err)
			time.Sleep(time.Second * time.Duration(p.ReloadTime))
			continue
		}
		for _, dbC := range dbContainers {
			var containerStillExist = false
			for _, c := range containers {
				if dbC.CID == c.ID {
					containerStillExist = true
					// Check if container started or stopped
					c1, ok1 := containers[dbC.CID]
					c2, ok2 := oldContainers[dbC.CID]
					if ok1 && ok2 && (c1.Running != c2.Running) {
						var event dguard.Event
						var eventSeverity int
						var eventType int
						if c1.Running {
							eventSeverity = dguard.EventNotice
							eventType = dguard.EventContainerStarted
						} else {
							eventSeverity = dguard.EventCritical
							eventType = dguard.EventContainerStopped
						}
						event = dguard.Event{
							eventSeverity,
							eventType,
							dbC.Hostname + " (" + dbC.CID + ")",
							p.Name,
							""}
						Alert(event)
					}
				}
			}
			if !containerStillExist {
				var event dguard.Event = dguard.Event{
					dguard.EventNotice,
					dguard.EventContainerRemoved,
					dbC.Hostname + " (" + dbC.CID + ")",
					p.Name,
					""}

				dbC.Delete()

				Alert(event)
			}
		}

		// Add containers and stats in DB
		for _, c := range containers {
			var id int64
			var tmpContainer Container

			// Add containers in DB
			tmpContainer, err = GetContainerByCID(c.ID)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					var event dguard.Event

					sqlContainer := Container{0, c.ID, probeID, c.Hostname, c.Image, c.IPAddress, c.MacAddress}

					event = dguard.Event{
						dguard.EventNotice,
						dguard.EventContainerCreated,
						sqlContainer.Hostname + " (" + sqlContainer.CID + ")",
						p.Name,
						"Image: " + sqlContainer.Image}
					id, err = sqlContainer.Insert()

					Alert(event)
					if err != nil {
						l.Error("MonitorProbe ("+p.Name+"): container insert:", err)
						continue
					}
				} else {
					l.Error("MonitorProbe ("+p.Name+"): GetContainerById:", err)
					continue
				}
			} else {
				id = int64(tmpContainer.ID)
			}

			// Add stats in DB
			sqlStat := Stat{int(id),
				int64(c.Time),
				uint64(c.SizeRootFs),
				uint64(c.SizeRw),
				uint64(c.MemoryUsed),
				uint64(c.NetBandwithRX),
				uint64(c.NetBandwithTX),
				uint64(c.CPUUsage),
				c.Running}
			err = sqlStat.Insert()
			if err != nil {
				l.Error("MonitorProbe ("+p.Name+"): stat insert:", err)
				continue
			}
		}

		time.Sleep(time.Second * time.Duration(p.ReloadTime))
	}
}
