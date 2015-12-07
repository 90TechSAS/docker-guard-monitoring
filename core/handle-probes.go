package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	dguard "github.com/90TechSAS/libgo-docker-guard"
	"github.com/gorilla/mux"
)

/*
	Return simplified probes array
*/
func HTTPHandlerProbes(w http.ResponseWriter, r *http.Request) {
	var returnStr string                 // HTTP Response body
	var returnProbes []dguard.ProbeInfos // Returned probes
	var populate string                  // HTTP GET parameter
	var err error                        // Error handling

	// Check if populate is true
	populate = r.URL.Query().Get("populate")

	// Get probes
	returnProbes = GetProbesInfos()

	// If populate == true, insert containers
	if populate == "true" {
		for i, probe := range returnProbes {
			returnProbes[i].Containers, err = GetSimpleContainersByProbe(probe.Name)
			if err != nil {
				l.Error("HTTPHandlerProbes: Failed to get list of containers")
				http.Error(w, http.StatusText(500), 500)
				return
			}
		}
	}

	// probes => json
	tmpJSON, err := json.Marshal(returnProbes)
	if err != nil {
		l.Error("HTTPHandlerProbes: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJSON)
	if returnStr == "null" {
		returnStr = "[]"
	}

	w.Header().Set("Content-Type", "application/json")
	AddCORS(w)
	fmt.Fprint(w, returnStr)
}

/*
	Return one probe
*/
func HTTPHandlerProbesName(w http.ResponseWriter, r *http.Request) {
	var returnStr string           // HTTP Response body
	var muxVars = mux.Vars(r)      // Mux Vars
	var probes []dguard.ProbeInfos // Probes
	var err error                  // Error handling

	// Get probe name
	probeNameVar, ok := muxVars["name"]
	if !ok {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Get probes
	probes = GetProbesInfos()

	// Search probe
	for _, p := range probes {
		if p.Name == probeNameVar {
			var sContainers []dguard.SimpleContainer

			// Get list of containers
			sContainers, err = GetSimpleContainersByProbe(p.Name)
			if err != nil {
				if err.Error() == "Not found" {
					http.Error(w, http.StatusText(404), 404)
					return
				}
				l.Error("HTTPHandlerProbesName: Failed to get list of containers")
				http.Error(w, http.StatusText(500), 500)
				return
			}

			// Add containers to probe infos
			p.Containers = sContainers

			// probes => json
			tmpJSON, err := json.Marshal(p)
			if err != nil {
				l.Error("HTTPHandlerProbesName: Failed to marshal struct")
				http.Error(w, http.StatusText(500), 500)
				return
			}

			// Add json to the returned string
			returnStr = string(tmpJSON)
			w.Header().Set("Content-Type", "application/json")
			AddCORS(w)
			fmt.Fprint(w, returnStr)
			return
		}

	}

	http.Error(w, http.StatusText(404), 404)
	return

	if returnStr == "null" {
		returnStr = "[]"
	}
}
