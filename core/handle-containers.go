package core

import (
	"encoding/json"
	"fmt"
	"net/http"

	dguard "github.com/90TechSAS/libgo-docker-guard"

	"../utils"
)

/*
	Return containers infos
*/
func HTTPHandlerContainers(w http.ResponseWriter, r *http.Request) {
	var returnStr string                      // HTTP Response body
	var tmpContainers []Container             // Temporary container list
	var returnedContainers []dguard.Container // Returned container list
	var err error                             // Error handling
	var probeid string                        // HTTP GET parameter(s)
	var probeidInt int                        // HTTP GET parameter(s) to int

	// Get HTTP query
	query := r.URL.Query()

	// Get HTTP GET parameters
	probeid = query.Get("probeid")

	// Get probe ID
	probeidInt, err = utils.S2I(probeid)
	if probeid == "" || err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	// Get containers by probe ID
	tmpContainers, err = GetContainersBy("probeid", probeidInt)
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

/*
	Return container infos
*/
func HTTPHandlerContainerID(w http.ResponseWriter, r *http.Request) {
	var returnStr string                   // HTTP Response body
	var returnedContainer dguard.Container // Returned container
	// var err error                          // Error handling

	// returnedContainer => json
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
	var returnedContainers []dguard.Container // Returned container list
	// var err error                             // Error handling

	// returnedContainers => json
	tmpJson, _ := json.Marshal(returnedContainers)

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
