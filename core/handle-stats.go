package core

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	var returnStr string     // HTTP Response body
	var err error            // Error handling
	var returnedStats []Stat // Returned stats

	http.Error(w, http.StatusText(501), 501) // Not implemented
	return

	// returnedStats => json
	tmpJson, err := json.Marshal(returnedStats)
	if err != nil {
		l.Error("HTTPHandlerStatsProbeID: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
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
	var returnStr string     // HTTP Response body
	var err error            // Error handling
	var returnedStats []Stat // Returned stats

	http.Error(w, http.StatusText(501), 501) // Not implemented
	return

	// returnedStats => json
	tmpJson, err := json.Marshal(returnedStats)
	if err != nil {
		l.Error("HTTPHandlerStatsCID: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
