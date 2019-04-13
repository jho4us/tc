package handlers

import (
	"github.com/gorilla/mux"

	"net/http"
	"sync/atomic"
)

// Router register extra kubernetes routes.
func MakeHandler(pathPrefix, commit, release string) http.Handler {
	isReady := &atomic.Value{}
	isReady.Store(true)

	r := mux.NewRouter()
	r.HandleFunc(pathPrefix+"/healthz", healthz)
	r.HandleFunc(pathPrefix+"/readyz", readyz(isReady))
	return r
}
