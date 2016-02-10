package main

import (
	"log"
	"net/http"
	"time"

	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/handlers"
	"github.com/deis/workflow-manager/jobs"
	"github.com/gorilla/mux"
)

func main() {
	// we want to do the following jobs according to our remote API interval:
	// 1. get latest stable deis component versions
	// 2. send diagnostic data, if appropriate
	toDo := []jobs.Periodic{jobs.GetLatestVersionData{}, jobs.SendVersions{}}
	ch := jobs.DoPeriodic(toDo, time.Duration(config.Spec.Polling))
	defer close(ch)
	// Get a new router, with handler functions
	r := handlers.RegisterRoutes(mux.NewRouter())
	// Bind to a port and pass our router in
	err := http.ListenAndServe(":"+config.Spec.Port, r)
	if err != nil {
		close(ch)
		log.Println("Unable to open up TLS listener")
		log.Fatal("ListenAndServe: ", err)
	}
}
