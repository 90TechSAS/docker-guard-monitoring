package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func HTTPHandlerProbes(w http.ResponseWriter, r *http.Request) {
	var returnStr string // HTTP Response body

	// probes => json
	tmpJson, _ := json.Marshal(ProbesID)

	// Add json to the returned string
	returnStr = string(tmpJson)

	fmt.Fprint(w, returnStr)
}
