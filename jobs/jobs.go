package jobs

import (
	"log"
	"time"

	"github.com/arschles/kubeapp/api/rc"
	"github.com/deis/workflow-manager/config"
	"github.com/deis/workflow-manager/data"
	apiclient "github.com/deis/workflow-manager/pkg/swagger/client"
	"github.com/deis/workflow-manager/pkg/swagger/client/operations"
)

// Periodic is an interface for managing periodic job invocation
type Periodic interface {
	// Do begins the periodic job. It starts the first execution of the job, and then is
	// repsonsible for executing it every Frequency() thereafter
	Do() error
	Frequency() time.Duration
}

// SendVersions fulfills the Periodic interface
type sendVersions struct {
	secretGetterCreator data.KubeSecretGetterCreator
	rcLister            rc.Lister
	apiClient           *apiclient.WorkflowManager
	availableVersions   data.AvailableVersions
	frequency           time.Duration
}

// NewSendVersionsPeriodic creates a new SendVersions using sgc and rcl as the the secret getter / creator and replication controller lister implementations (respectively)
func NewSendVersionsPeriodic(
	apiClient *apiclient.WorkflowManager,
	sgc data.KubeSecretGetterCreator,
	rcl rc.Lister,
	availableVersions data.AvailableVersions,
	frequency time.Duration,
) Periodic {
	return &sendVersions{
		secretGetterCreator: sgc,
		rcLister:            rcl,
		apiClient:           apiClient,
		availableVersions:   availableVersions,
		frequency:           frequency,
	}
}

// Do is the Periodic interface implementation
func (s sendVersions) Do() error {
	if config.Spec.CheckVersions {
		err := sendVersionsImpl(s.apiClient, s.secretGetterCreator, s.rcLister, s.availableVersions)
		if err != nil {
			return err
		}
	}
	return nil
}

// Frequency is the Periodic interface implementation
func (s sendVersions) Frequency() time.Duration {
	return s.frequency
}

type getLatestVersionData struct {
	vsns                  data.AvailableVersions
	installedData         data.InstalledData
	clusterID             data.ClusterID
	availableComponentVsn data.AvailableComponentVersion
	sgc                   data.KubeSecretGetterCreator
	frequency             time.Duration
}

// NewGetLatestVersionDataPeriodic creates a new periodic implementation that gets latest version data. It uses sgc and rcl as the secret getter/creator and replication controller lister implementations (respectively)
func NewGetLatestVersionDataPeriodic(
	sgc data.KubeSecretGetterCreator,
	rcl rc.Lister,
	installedData data.InstalledData,
	clusterID data.ClusterID,
	availVsn data.AvailableVersions,
	availCompVsn data.AvailableComponentVersion,
	frequency time.Duration,
) Periodic {

	return &getLatestVersionData{
		vsns:                  availVsn,
		installedData:         installedData,
		clusterID:             clusterID,
		availableComponentVsn: availCompVsn,
		sgc:       sgc,
		frequency: frequency,
	}
}

// Do is the Periodic interface implementation
func (u *getLatestVersionData) Do() error {
	cluster, err := data.GetCluster(u.installedData, u.clusterID, u.availableComponentVsn, u.sgc)
	if err != nil {
		return err
	}
	if _, err := u.vsns.Refresh(cluster); err != nil {
		return err
	}
	return nil
}

// Frequency is the Periodic interface implementation
func (u getLatestVersionData) Frequency() time.Duration {
	return u.frequency
}

// DoPeriodic calls p.Do() every p.Frequency() on each element p in pSlice. For each p in pSlice,
// a new goroutine is started, and the returned channel can be closed to stop all of the goroutines
func DoPeriodic(pSlice []Periodic) chan<- struct{} {
	doneCh := make(chan struct{})
	for _, p := range pSlice {
		go func(p Periodic) {
			ticker := time.NewTicker(p.Frequency())
			for {
				select {
				case <-ticker.C:
					err := p.Do()
					if err != nil {
						log.Printf("periodic job ran and returned error (%s)", err)
					}
				case <-doneCh:
					ticker.Stop()
					return
				}
			}
		}(p)
	}
	return doneCh
}

//  sendVersions sends cluster version data
func sendVersionsImpl(
	apiClient *apiclient.WorkflowManager,
	secretGetterCreator data.KubeSecretGetterCreator,
	rcLister rc.Lister,
	availableVersions data.AvailableVersions,
) error {
	cluster, err := data.GetCluster(
		data.NewInstalledDeisData(rcLister),
		data.NewClusterIDFromPersistentStorage(secretGetterCreator),
		data.NewLatestReleasedComponent(secretGetterCreator, rcLister, availableVersions),
		secretGetterCreator,
	)
	if err != nil {
		log.Println("error getting installed components data")
		return err
	}

	_, err = apiClient.Operations.CreateClusterDetails(&operations.CreateClusterDetailsParams{Body: &cluster})
	if err != nil {
		log.Println("error sending diagnostic data")
		return err
	}
	return nil
}
