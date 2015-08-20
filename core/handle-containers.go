package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	dguard "github.com/90TechSAS/libgo-docker-guard"
	"github.com/gorilla/mux"

	"../utils"
)

/*
	Return containers infos
*/
func HTTPHandlerContainers(w http.ResponseWriter, r *http.Request) {
	var returnStr string               // HTTP Response body
	var returnedContainers []Container // Returned container
	var err error                      // Error handling

	http.Error(w, http.StatusText(501), 501) // Not implemented

	// returnedContainers => json
	tmpJson, err := json.Marshal(returnedContainers)
	if err != nil {
		l.Error("HTTPHandlerContainers: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}

/*
	Return container infos
*/
func HTTPHandlerContainerID(w http.ResponseWriter, r *http.Request) {
	var returnStr string            // HTTP Response body
	var returnedContainer Container // Returned container
	var muxVars = mux.Vars(r)       // Mux Vars
	var err error                   // Error handling

	// Get container ID
	ContainerIDVar := muxVars["id"]
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Get container
	containers, err := GetContainersBy("containerid", ContainerIDVar)
	if err != nil || len(containers) > 1 {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	if err != nil || len(containers) == 0 {
		http.Error(w, http.StatusText(404), 404)
		return
	}

	// returnedContainer => json
	returnedContainer = containers[0]
	tmpJson, err := json.Marshal(returnedContainer)
	if err != nil {
		l.Error("HTTPHandlerContainerID: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

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
	probeIDVar, err := utils.S2I(muxVars["id"])
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Get containers by probe ID
	returnedContainers, err = GetContainersBy("probeid", probeIDVar)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// returnedContainers => json
	tmpJson, err := json.Marshal(returnedContainers)
	if err != nil {
		l.Error("HTTPHandlerContainersProbeID: Failed to marshal struct")
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
