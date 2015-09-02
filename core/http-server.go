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
	// If HTTP Method is OPTIONS, don't check basic auth
	if r.Method == "OPTIONS" {
		return false
	}

	// Check basic auth
	u, p, ok := r.BasicAuth()
	if ok == true && u == DGConfig.DockerGuard.API.APILogin && p == DGConfig.DockerGuard.API.APIPassword {
		l.Debug("Auth OK from", r.RemoteAddr)
		return true
	}

	// If auth is not ok, check why
	if !ok {
		l.Warn("Failed auth from", r.RemoteAddr, ", basic auth is not ok")
	} else {
		l.Warn("Failed auth from", r.RemoteAddr, ", credentials are not ok")
	}

	return false
}

/*
	Add CORS to headers
*/
func AddCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	// If HTTP Method is OPTIONS, return CORS
	if r.Method == "OPTIONS" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization,DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "1728000")
		w.Header().Set("Content-Type", "text/plain charset=UTF-8")
		w.Header().Set("Content-Length", "0")
		l.Info("OPTIONS")
		return
	}

	// If header "Authorization" is not null and invalid, return 403
	if r.Header.Get("Authorization") != "" {
		u, p, ok := r.BasicAuth()
		if ok == true && u == DGConfig.DockerGuard.API.APILogin && p == DGConfig.DockerGuard.API.APIPassword {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		http.Error(w, http.StatusText(403), 403)
	} else {
		http.Error(w, http.StatusText(404), 404)
	}
}

/*
	Run HTTP Server
*/
func HTTPServer() {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(NotFoundHandler)
	r1 := r.MatcherFunc(HTTPURILogger).MatcherFunc(HTTPSecureAPI).Subrouter()
	// r1 := r.MatcherFunc(HTTPURILogger).Subrouter()
	rGET := r1.Methods("GET").Subrouter()
	// rOPTIONS := r.MatcherFunc(HTTPURILogger).Methods("OPTIONS").Subrouter()

	rGET.HandleFunc("/containers", HTTPHandlerContainers)
	rGET.HandleFunc("/containers/{cid:[0-9a-z]+}", HTTPHandlerContainerCID)
	rGET.HandleFunc("/containers/probe/{name:[0-9a-zA-Z-_]+}", HTTPHandlerContainersProbeName)
	rGET.HandleFunc("/probes", HTTPHandlerProbes)
	rGET.HandleFunc("/probes/{name:[0-9a-zA-Z-_]+}", HTTPHandlerProbesName)
	rGET.HandleFunc("/stats", HTTPHandlerStats)
	rGET.HandleFunc("/stats/probe/{name:[0-9a-zA-Z-_]+}", HTTPHandlerStatsProbeName)
	rGET.HandleFunc("/stats/container/{cid:[0-9a-z]+}", HTTPHandlerStatsCID)
	http.Handle("/", r)

	http.ListenAndServe(DGConfig.DockerGuard.API.ListenInterface+":"+DGConfig.DockerGuard.API.ListenPort, r)
}
