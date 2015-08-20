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
func RunHTTPServer() {
	r := mux.NewRouter()
	r1 := r.MatcherFunc(HTTPURILogger).MatcherFunc(HTTPSecureAPI).Subrouter()
	r_GET := r1.Methods("GET").Subrouter()

	r_GET.HandleFunc("/containers", HTTPHandlerContainers)
	r_GET.HandleFunc("/containers/{cid:[0-9a-z]+}", HTTPHandlerContainerCID)
	r_GET.HandleFunc("/containers/probe/{id:[0-9]+}", HTTPHandlerContainersProbeID)
	r_GET.HandleFunc("/probes", HTTPHandlerProbes)
	r_GET.HandleFunc("/probes/{id:[0-9]+}", HTTPHandlerProbesID)
	r_GET.HandleFunc("/stats", HTTPHandlerStats)
	r_GET.HandleFunc("/stats/probe/{id:[0-9]+}", HTTPHandlerStatsProbeID)
	r_GET.HandleFunc("/stats/container/{cid:[0-9a-z]+}", HTTPHandlerStatsCID)
	http.Handle("/", r)

	http.ListenAndServe(DGConfig.DockerGuard.API.ListenInterface+":"+DGConfig.DockerGuard.API.ListenPort, r)
}
