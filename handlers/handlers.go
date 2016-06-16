package handlers

// handler echoes the HTTP request.
import (
	"encoding/json"
	"net/http"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	"github.com/deis/workflow-manager/pkg/swagger/client/operations"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

const (
	componentsRoute = "/components" // resource value for components route
	idRoute         = "/id"         // resource value for ID route
	doctorRoute     = "/doctor"
)

// RegisterRoutes attaches handler functions to routes
func RegisterRoutes(
	r *mux.Router,
	secretGetterCreator data.KubeSecretGetterCreator,
	rcLister rc.Lister,
	availableVersions data.AvailableVersions,
) *mux.Router {

	clusterID := data.NewClusterIDFromPersistentStorage(secretGetterCreator)
	r.Handle(componentsRoute, ComponentsHandler(
		data.NewInstalledDeisData(rcLister),
		clusterID,
		data.NewLatestReleasedComponent(secretGetterCreator, rcLister, availableVersions),
		secretGetterCreator,
	))
	r.Handle(idRoute, IDHandler(clusterID))
	r.Handle(doctorRoute, DoctorHandler(
		data.NewInstalledDeisData(rcLister),
		clusterID,
		data.NewLatestReleasedComponent(secretGetterCreator, rcLister, availableVersions),
		secretGetterCreator,
	)).Methods("POST")
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

// DoctorHandler route handler
func DoctorHandler(
	c data.InstalledData,
	i data.ClusterID,
	v data.AvailableComponentVersion,
	secretGetterCreator data.KubeSecretGetterCreator,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doctor, err := data.GetDoctorInfo(c, i, v, secretGetterCreator)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		apiClient, err := config.GetSwaggerClient(config.Spec.DoctorAPIURL)
		uid := uuid.NewV4().String()
		_, err = apiClient.Operations.PublishDoctorInfo(&operations.PublishDoctorInfoParams{Body: &doctor, UUID: uid})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writePlainText(uid, w)
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
