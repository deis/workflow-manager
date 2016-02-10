package handlers

// handler echoes the HTTP request.
import (
	"encoding/json"
	"net/http"

	"github.com/deis/workflow-manager/components"
	"github.com/deis/workflow-manager/data"
	"github.com/gorilla/mux"
)

const (
	componentsRoute = "/components" // resource value for components route
	idRoute         = "/id"         // resource value for ID route
)

// RegisterRoutes attaches handler functions to routes
func RegisterRoutes(r *mux.Router) *mux.Router {
	r.Handle(componentsRoute, ComponentsHandler(components.InstalledDeisData{}, data.ClusterIDFromPersistentStorage{}, components.LatestReleasedComponent{}))
	r.Handle(idRoute, IDHandler(data.ClusterIDFromPersistentStorage{}))
	return r
}

// ComponentsHandler route handler
func ComponentsHandler(c components.InstalledData, i data.ClusterID, v components.AvailableComponentVersion) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cluster, err := components.GetCluster(c, i, v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(cluster); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

// IDHandler route handler
func IDHandler(getter data.ClusterID) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := data.GetID(getter)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writePlainText(id, w)
	})
}

// writePlainText is a helper function for writing HTTP text data
func writePlainText(text string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(text))
}
