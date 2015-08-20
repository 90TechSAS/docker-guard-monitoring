package core

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	tmpJson, err := json.Marshal(returnProbes)
	if err != nil {
		l.Error("HTTPHandlerProbes: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)
	fmt.Fprint(w, returnStr)
}
