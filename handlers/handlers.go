package handlers

// handler echoes the HTTP request.
import (
	"encoding/json"
	"net/http"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	"github.com/deis/workflow-manager/k8s"
	apiclient "github.com/deis/workflow-manager/pkg/swagger/client"
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
	availVers data.AvailableVersions,
	k8sResources *k8s.ResourceInterfaceNamespaced,
) *mux.Router {

	clusterID := data.NewClusterIDFromPersistentStorage(k8sResources.Secrets())
	r.Handle(componentsRoute, ComponentsHandler(
		data.NewInstalledDeisData(k8sResources),
		clusterID,
		data.NewLatestReleasedComponent(k8sResources, availVers),
	))
	r.Handle(idRoute, IDHandler(clusterID))
	doctorAPIClient, _ := config.GetSwaggerClient(config.Spec.DoctorAPIURL)
	r.Handle(doctorRoute, DoctorHandler(
		data.NewInstalledDeisData(k8sResources),
		k8s.NewRunningK8sData(k8sResources),
		clusterID,
		data.NewLatestReleasedComponent(k8sResources, availVers),
		doctorAPIClient,
	)).Methods("POST")
	return r
}

// ComponentsHandler route handler
func ComponentsHandler(
	workflow data.InstalledData,
	clusterID data.ClusterID,
	availVers data.AvailableComponentVersion,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cluster, err := data.GetCluster(workflow, clusterID, availVers)
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
	workflow data.InstalledData,
	k8sData k8s.RunningK8sData,
	clusterID data.ClusterID,
	availVers data.AvailableComponentVersion,
	apiClient *apiclient.WorkflowManager,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		doctor, err := data.GetDoctorInfo(workflow, k8sData, clusterID, availVers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
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
