package core

import (
	"net/http"

	"github.com/gorilla/mux"
)

/*
	Log HTTP requests' URI
*/
func HTTPURILogger(r *http.Request, rm *mux.RouteMatch) bool {
	l.Verbose("Request URI:", r.RequestURI)
	return true
}

/*
	Secure API access
*/
func HTTPSecureAPI(r *http.Request, rm *mux.RouteMatch) bool {
	auth, ok := r.Header["Auth"]

	if ok && len(auth) == 1 && auth[0] == DGConfig.DockerGuard.API.APIPassword {
		l.Debug("Auth OK from", r.RemoteAddr)
		return true
	}

	l.Warn("Failed auth from", r.RemoteAddr)
	return false
}

/*
	Run HTTP Server
*/
func HTTPServer() {
	r := mux.NewRouter()
	r1 := r.MatcherFunc(HTTPURILogger).MatcherFunc(HTTPSecureAPI).Subrouter()
	rGET := r1.Methods("GET").Subrouter()

	rGET.HandleFunc("/containers", HTTPHandlerContainers)
	rGET.HandleFunc("/containers/{cid:[0-9a-z]+}", HTTPHandlerContainerCID)
	rGET.HandleFunc("/containers/probe/{id:[0-9]+}", HTTPHandlerContainersProbeID)
	rGET.HandleFunc("/probes", HTTPHandlerProbes)
	rGET.HandleFunc("/probes/{id:[0-9]+}", HTTPHandlerProbesID)
	rGET.HandleFunc("/stats", HTTPHandlerStats)
	rGET.HandleFunc("/stats/probe/{id:[0-9]+}", HTTPHandlerStatsProbeID)
	rGET.HandleFunc("/stats/container/{cid:[0-9a-z]+}", HTTPHandlerStatsCID)
	http.Handle("/", r)

	http.ListenAndServe(DGConfig.DockerGuard.API.ListenInterface+":"+DGConfig.DockerGuard.API.ListenPort, r)
}
