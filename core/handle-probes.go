package core

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"

	"../utils"
)

/*
	Simplified Probe struct
*/
type SProbe struct {
	ID   int    `yaml:"id"`
	Name string `yaml:"name"`
}

/*
	Return simplified probes array
*/
func HTTPHandlerProbes(w http.ResponseWriter, r *http.Request) {
	var returnStr string      // HTTP Response body
	var returnProbes []SProbe // Returned simplified probes
	var err error             // Error handling

	// Get simplified probes
	for key, probeID := range ProbesID {
		returnProbes = append(returnProbes, SProbe{probeID, key})
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
	fmt.Fprint(w, returnStr)
}

/*
	Return one simplified probe
*/
func HTTPHandlerProbesID(w http.ResponseWriter, r *http.Request) {
	var returnStr string      // HTTP Response body
	var returnProbe SProbe    // Returned simplified probe
	var err error             // Error handling
	var muxVars = mux.Vars(r) // Mux Vars
	var probeFound = false

	// Get simplified probe
	for key, probeID := range ProbesID {
		probeIDVar, err := utils.S2I(muxVars["id"])
		if err != nil {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		if probeID == probeIDVar {
			returnProbe = SProbe{probeID, key}
			probeFound = true
			break
		}
	}

	// Check if found
	if !probeFound {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// Get simplified probe
	// probe => json
	tmpJSON, err := json.Marshal(returnProbe)
	if err != nil {
		l.Error("HTTPHandlerProbesID: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJSON)
	fmt.Fprint(w, returnStr)
}
