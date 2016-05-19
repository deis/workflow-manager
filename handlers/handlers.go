package handlers

// handler echoes the HTTP request.
import (
	"encoding/json"
	"net/http"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/data"
	"github.com/gorilla/mux"
)

const (
	componentsRoute = "/components" // resource value for components route
	idRoute         = "/id"         // resource value for ID route
)

// RegisterRoutes attaches handler functions to routes
func RegisterRoutes(r *mux.Router, secretGetterCreator data.KubeSecretGetterCreator, rcLister rc.Lister) *mux.Router {
	clusterID := data.NewClusterIDFromPersistentStorage(secretGetterCreator)
	r.Handle(componentsRoute, ComponentsHandler(
		data.InstalledDeisData{},
		clusterID,
		data.NewLatestReleasedComponent(secretGetterCreator, rcLister),
		secretGetterCreator,
	))
	r.Handle(idRoute, IDHandler(clusterID))
	return r
}

// ComponentsHandler route handler
func ComponentsHandler(
	c data.InstalledData,
	i data.ClusterID,
	v data.AvailableComponentVersion,
	secretGetterCreator data.KubeSecretGetterCreator,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cluster, err := data.GetCluster(c, i, v, secretGetterCreator)
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
