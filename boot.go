package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	"github.com/deis/workflow-manager/handlers"
	"github.com/deis/workflow-manager/jobs"
	"github.com/gorilla/mux"
	kcl "k8s.io/kubernetes/pkg/client/unversioned"
)

func main() {
	kubeClient, err := kcl.NewInCluster()
	if err != nil {
		log.Fatalf("Error creating new Kubernetes client (%s)", err)
	}
	apiClient := config.GetSwaggerClient(config.Spec.VersionsAPIURL)
	secretInterface := kubeClient.Secrets(config.Spec.DeisNamespace)
	rcInterface := kubeClient.ReplicationControllers(config.Spec.DeisNamespace)
	clusterID := data.NewClusterIDFromPersistentStorage(secretInterface)
	installedDeisData := data.NewInstalledDeisData(rcInterface)
	availableVersion := data.NewAvailableVersionsFromAPI(
		apiClient,
		config.Spec.VersionsAPIURL,
		secretInterface,
		rcInterface,
	)
	availableComponentVersion := data.NewLatestReleasedComponent(secretInterface, rcInterface, availableVersion)

	// we want to do the following jobs according to our remote API interval:
	// 1. get latest stable deis component versions
	// 2. send diagnostic data, if appropriate
	glvdPeriodic := jobs.NewGetLatestVersionDataPeriodic(
		secretInterface,
		rcInterface,
		installedDeisData,
		clusterID,
		availableVersion,
		availableComponentVersion,
	)
	svPeriodic := jobs.NewSendVersionsPeriodic(apiClient, secretInterface, rcInterface, availableVersion)
	toDo := []jobs.Periodic{glvdPeriodic, svPeriodic}
	pollDur := time.Duration(config.Spec.Polling) * time.Second
	log.Printf("Starting periodic jobs at interval %s", pollDur)
	ch := jobs.DoPeriodic(toDo, time.Duration(config.Spec.Polling)*time.Second)
	defer close(ch)

	// Get a new router, with handler functions
	r := handlers.RegisterRoutes(mux.NewRouter(), secretInterface, rcInterface, availableVersion)
	// Bind to a port and pass our router in
	hostStr := fmt.Sprintf(":%s", config.Spec.Port)
	log.Printf("Serving on %s", hostStr)
	if err := http.ListenAndServe(hostStr, r); err != nil {
		close(ch)
		log.Println("Unable to open up TLS listener")
		log.Fatal("ListenAndServe: ", err)
	}
}
