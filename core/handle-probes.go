package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	dguard "github.com/90TechSAS/libgo-docker-guard"
)

/*
	Return simplified probes array
*/
func HTTPHandlerProbes(w http.ResponseWriter, r *http.Request) {
	var returnStr string                 // HTTP Response body
	var returnProbes []dguard.ProbeInfos // Returned probes
	var err error                        // Error handling

	// Get probes
	returnProbes = GetProbesInfos()

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
func HTTPHandlerProbesID(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(501), 501) // Not implemented
	return
}
