package core

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

/*
	Stat populated
*/
type StatPopulated struct {
	Container     Container
	Time          time.Time
	SizeRootFs    uint64
	SizeRw        uint64
	SizeMemory    uint64
	NetBandwithRX uint64
	NetBandwithTX uint64
	CPUUsage      uint64
	Running       bool
}

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
	tmpJSON, err := json.Marshal(returnedStats)
	if err != nil {
		l.Error("HTTPHandlerStats: Failed to marshal struct")
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
	Return containers stats by probe name
*/
func HTTPHandlerStatsProbeName(w http.ResponseWriter, r *http.Request) {
	var returnStr string      // HTTP Response body
	var muxVars = mux.Vars(r) // Mux Vars
	var tmpJSON []byte        // Temporary JSON
	var options Options       // Options
	var err error             // Error handling

	options = GetOptions(r)

	// Get mux Vars
	probeNameVar := muxVars["name"]

	// Check if populate is wanted
	populate := r.URL.Query().Get("populate")
	if populate == "true" {
		var returnedStats []StatPopulated // Returned stats
		returnedStats, err = GetStatsPByContainerProbeID(probeNameVar, options)
		if err != nil {
			if strings.Contains(err.Error(), "Not found") {
				l.Error("HTTPHandlerStatsProbeName: Failed to get stats:", err)
				http.Error(w, http.StatusText(404), 404)
				return
			}
			l.Error("HTTPHandlerStatsProbeName: Failed to get stats:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// returnedStats => json
		tmpJSON, err = json.Marshal(returnedStats)
		if err != nil {
			l.Error("HTTPHandlerStatsProbeName: Failed to marshal struct:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	} else if populate == "false" || populate == "" {
		var returnedStats []Stat // Returned stats
		returnedStats, err = GetStatsByContainerProbeID(probeNameVar, options)
		if err != nil {
			if strings.Contains(err.Error(), "Not found") {
				l.Error("HTTPHandlerStatsProbeName: Failed to get stats:", err)
				http.Error(w, http.StatusText(404), 404)
				return
			}
			l.Error("HTTPHandlerStatsProbeName: Failed to get stats:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}

		// returnedStats => json
		tmpJSON, err = json.Marshal(returnedStats)
		if err != nil {
			l.Error("HTTPHandlerStatsProbeName: Failed to marshal struct:", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	} else {
		http.Error(w, http.StatusText(400), 400)
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
	Return containers stats by container ID
*/
func HTTPHandlerStatsCID(w http.ResponseWriter, r *http.Request) {
	var returnStr string      // HTTP Response body
	var returnedStats []Stat  // Returned stats
	var muxVars = mux.Vars(r) // Mux Vars
	var err error             // Error handling
	var options Options       // Options

	options = GetOptions(r)

	// Set default Limit
	if options.Limit == -1 {
		options.Limit = 20
	}

	// Get mux Vars
	containerCIDVar := muxVars["cid"]

	returnedStats, err = GetStatsByContainerCID(containerCIDVar, options)
	if err != nil {
		l.Error("HTTPHandlerStatsCID: Failed to get stats:", err)
		if strings.Contains(err.Error(), "Not found") {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// returnedStats => json
	tmpJSON, err := json.Marshal(returnedStats)
	if err != nil {
		l.Error("HTTPHandlerStatsCID: Failed to marshal struct:", err)
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
