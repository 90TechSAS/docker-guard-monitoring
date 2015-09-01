package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

/*
	Return containers infos
*/
func HTTPHandlerContainers(w http.ResponseWriter, r *http.Request) {
	var returnStr string               // HTTP Response body
	var returnedContainers []Container // Returned container
	var err error                      // Error handling

	http.Error(w, http.StatusText(501), 501) // Not implemented
	return

	// returnedContainers => json
	tmpJSON, err := json.Marshal(returnedContainers)
	if err != nil {
		l.Error("HTTPHandlerContainers: Failed to marshal struct")
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
	Return container infos
*/
func HTTPHandlerContainerCID(w http.ResponseWriter, r *http.Request) {
	var returnStr string            // HTTP Response body
	var returnedContainer Container // Returned container
	var muxVars = mux.Vars(r)       // Mux Vars
	var err error                   // Error handling

	// Get container ID
	ContainerIDVar := muxVars["cid"]
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Get container
	returnedContainer, err = GetContainerByCID(ContainerIDVar)
	if err != nil {
		if err.Error() == "Not Found" {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// returnedContainer => json
	tmpJSON, err := json.Marshal(returnedContainer)
	if err != nil {
		l.Error("HTTPHandlerContainerID: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJSON)
	if returnStr == "null" {
		returnStr = "{}"
	}

	w.Header().Set("Content-Type", "application/json")
	AddCORS(w)
	fmt.Fprint(w, returnStr)
}

/*
	Return probe's containers infos
*/
func HTTPHandlerContainersProbeID(w http.ResponseWriter, r *http.Request) {
	var returnStr string               // HTTP Response body
	var returnedContainers []Container // Returned container list
	var muxVars = mux.Vars(r)          // Mux Vars
	var err error                      // Error handling

	// Get probe ID
	probeIDVar := muxVars["id"]

	// Get containers by probe ID
	returnedContainers, err = GetContainersByProbe(probeIDVar)
	if err != nil {
		if err.Error() == "Not Found" {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// returnedContainers => json
	tmpJSON, err := json.Marshal(returnedContainers)
	if err != nil {
		l.Error("HTTPHandlerContainersProbeID: Failed to marshal struct")
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
