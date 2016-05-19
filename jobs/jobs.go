package jobs

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
)

// Periodic is an interface for managing periodic job invocation
type Periodic interface {
	// will have a Do method to begin execution
	Do() error
}

// SendVersions fulfills the Periodic interface
type SendVersions struct {
	secretGetterCreator data.KubeSecretGetterCreator
	rcLister            rc.Lister
}

// NewSendVersions creates a new SendVersions using sgc and rcl as the the secret getter / creator and replication controller lister implementations (respectively)
func NewSendVersions(sgc data.KubeSecretGetterCreator, rcl rc.Lister) *SendVersions {
	return &SendVersions{secretGetterCreator: sgc, rcLister: rcl}
}

// Do method of SendVersions
func (s SendVersions) Do() error {
	if config.Spec.CheckVersions {
		err := sendVersions(s.secretGetterCreator, s.rcLister)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLatestVersionData fulfills the Periodic interface
type GetLatestVersionData struct {
	secretGetterCreator data.KubeSecretGetterCreator
	rcLister            rc.Lister
}

// NewGetLatestVersionData creates a new GetLatestVersionData using sgc and rcl as the secret getter/creator and replication controller lister implementations (respectively)
func NewGetLatestVersionData(sgc data.KubeSecretGetterCreator, rcl rc.Lister) *GetLatestVersionData {
	return &GetLatestVersionData{secretGetterCreator: sgc, rcLister: rcl}
}

// Do method of GetLatestVersionData
func (u GetLatestVersionData) Do() error {
	dataSource := data.NewAvailableVersionsFromAPI("", u.secretGetterCreator, u.rcLister)
	_, err := dataSource.Refresh()
	if err != nil {
		return err
	}
	return nil
}

// DoPeriodic is a function for running jobs at a fixed interval
func DoPeriodic(p []Periodic, interval time.Duration) chan struct{} {
	ch := make(chan struct{})
	// schedule later job runs at a regular, periodic interval
	ticker := time.NewTicker(interval * time.Second)
	go func() {
		// run the period jobs once at invocation time
		runJobs(p)
		for {
			select {
			case <-ticker.C:
				runJobs(p)
			case <-ch:
				ticker.Stop()
				return
			}
		}
	}()
	return ch
}

// runJobs is a helper function to run a list of jobs
func runJobs(p []Periodic) {
	for _, job := range p {
		err := job.Do()
		if err != nil {
			log.Println("periodic job ran and returned error:")
			log.Print(err)
		}
	}
}

//  sendVersions sends cluster version data
func sendVersions(secretGetterCreator data.KubeSecretGetterCreator, rcLister rc.Lister) error {
	var clustersRoute = "/" + config.Spec.APIVersion + "/clusters/"
	cluster, err := data.GetCluster(
		data.InstalledDeisData{},
		data.NewClusterIDFromPersistentStorage(secretGetterCreator),
		data.NewLatestReleasedComponent(secretGetterCreator, rcLister),
		secretGetterCreator,
	)
	if err != nil {
		log.Println("error getting installed components data")
		return err
	}
	url := config.Spec.VersionsAPIURL + clustersRoute + cluster.ID
	js, err := json.Marshal(cluster)
	if err != nil {
		log.Println("error making a JSON representation of cluster data")
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	if err != nil {
		log.Println("error constructing POST request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := getTLSClient().Do(req)
	if err != nil {
		log.Println("error sending diagnostic data")
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getTLSClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	return &http.Client{Transport: tr}
}
