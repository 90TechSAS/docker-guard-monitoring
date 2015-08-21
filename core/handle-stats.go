package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

/*
	Return containers stats
*/
func HTTPHandlerStats(w http.ResponseWriter, r *http.Request) {
	var returnStr string     // HTTP Response body
	var err error            // Error handling
	var returnedStats []Stat // Returned stats

	http.Error(w, http.StatusText(501), 501) // Not implemented
	return

	// returnedStats => json
	tmpJson, err := json.Marshal(returnedStats)
	if err != nil {
		l.Error("HTTPHandlerStats: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}

/*
	Return containers stats by probe ID
*/
func HTTPHandlerStatsProbeID(w http.ResponseWriter, r *http.Request) {
	var returnStr string      // HTTP Response body
	var muxVars = mux.Vars(r) // Mux Vars
	var tmpJson []byte        // Temporary JSON
	var options Options       // Options
	var err error             // Error handling

	options = GetOptions(r)

	// Get mux Vars
	probeIDVar := muxVars["id"]

	// Check if populate is wanted
	populate := r.URL.Query().Get("populate")
	if populate == "true" {
		var returnedStats []StatPopulated // Returned stats
		returnedStats, err = GetStatsPByContainerProbeID(probeIDVar, options)
		if err != nil {
			l.Error("HTTPHandlerStatsProbeID: Failed to get stats:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// returnedStats => json
		tmpJson, err = json.Marshal(returnedStats)
		if err != nil {
			l.Error("HTTPHandlerStatsProbeID: Failed to marshal struct:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	} else if populate == "false" || populate == "" {
		var returnedStats []Stat // Returned stats
		returnedStats, err = GetStatsByContainerProbeID(probeIDVar, options)
		if err != nil {
			l.Error("HTTPHandlerStatsProbeID: Failed to get stats:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// returnedStats => json
		tmpJson, err = json.Marshal(returnedStats)
		if err != nil {
			l.Error("HTTPHandlerStatsProbeID: Failed to marshal struct:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	} else {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}

/*
	Return containers stats by container ID
*/
func HTTPHandlerStatsCID(w http.ResponseWriter, r *http.Request) {
	var returnStr string      // HTTP Response body
	var returnedStats []Stat  // Returned stats
	var muxVars = mux.Vars(r) // Mux Vars
	var err error             // Error handling
	var options Options       // Options

	options = GetOptions(r)

	// Get mux Vars
	containerCIDVar := muxVars["cid"]

	returnedStats, err = GetStatsByContainerCID(containerCIDVar, options)
	if err != nil {
		l.Error("HTTPHandlerStatsCID: Failed to get stats:", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// returnedStats => json
	tmpJson, err := json.Marshal(returnedStats)
	if err != nil {
		l.Error("HTTPHandlerStatsCID: Failed to marshal struct:", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
