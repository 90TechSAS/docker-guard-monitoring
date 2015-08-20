package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

/*
	Return containers infos
*/
func HTTPHandlerStats(w http.ResponseWriter, r *http.Request) {
	var returnStr string // HTTP Response body
	// var err error            // Error handling
	var returnedStats []Stat // Returned stats

	// returnedContainers => json
	tmpJson, _ := json.Marshal(returnedStats)

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
