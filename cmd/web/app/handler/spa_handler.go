package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
)

// SpaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type SpaHandler struct {
	StaticPath string
	IndexPath  string
}

// HTTPHandlerParams contains parameters that must be passed to NewHTTPHandler.
type HTTPHandlerParams struct {

	// BasePath will be used as a prefix for the endpoints, e.g. "/api"
	BasePath string
}

// NewHTTPHandler creates new HTTPHandler.
func NewHTTPHandler(params HTTPHandlerParams) *HTTPHandler {
	handler := &HTTPHandler{params: params}

	return handler
}

// HTTPHandler holds the params for a HTTPHandler
type HTTPHandler struct {
	params HTTPHandlerParams
}

// RegisterRoutes registers routes for this handler on the given router
func (sH *SpaHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/", sH.ServeHTTP).Methods(http.MethodGet)
}

func (h SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request
		// and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.StaticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.StaticPath, h.IndexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.StaticPath)).ServeHTTP(w, r)
}
