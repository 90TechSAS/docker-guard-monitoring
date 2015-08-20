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
	var returnStr string                      // HTTP Response body
	var returnedContainers []dguard.Container // Returned container
	// var err error                          // Error handling

	http.Error(w, http.StatusText(501), 501) // Not implemented

	// returnedContainers => json
	tmpJson, _ := json.Marshal(returnedContainers)

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
	tmpJson, _ := json.Marshal(returnedContainer)

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}

/*
	Return probe's containers infos
*/
func HTTPHandlerContainersProbeID(w http.ResponseWriter, r *http.Request) {
	var returnStr string                      // HTTP Response body
	var tmpContainers []Container             // Temporary container list
	var returnedContainers []dguard.Container // Returned container list
	var muxVars = mux.Vars(r)                 // Mux Vars
	var err error                             // Error handling

	probeIDVar, err := utils.S2I(muxVars["id"])
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Get containers by probe ID
	tmpContainers, err = GetContainersBy("probeid", probeIDVar)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	// Get containers last stats
	for _, c := range tmpContainers {
		var tmpC dguard.Container
		var tmpStat Stat

		tmpStat, err = c.GetLastStat()
		if err != nil {
			l.Critical("err:", err)
		}

		tmpC.ID = c.CID
		tmpC.Hostname = c.Hostname
		tmpC.Image = c.Image
		tmpC.IPAddress = c.IPAddress
		tmpC.MacAddress = c.MacAddress
		tmpC.SizeRootFs = float64(tmpStat.SizeRootFs)
		tmpC.SizeRw = float64(tmpStat.SizeRw)
		tmpC.MemoryUsed = float64(tmpStat.SizeMemory)
		tmpC.Running = tmpStat.Running
		tmpC.Time = float64(tmpStat.Time)

		returnedContainers = append(returnedContainers, tmpC)
	}

	// returnedContainers => json
	tmpJson, _ := json.Marshal(returnedContainers)

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
