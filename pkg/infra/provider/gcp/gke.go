package gcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"

	container "cloud.google.com/go/container/apiv1"
	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/infra/provider/kubernetes/helm"
	"github.com/puppetlabs/nebula/pkg/state"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	defaultMachineType      = "f1-micro"
	defaultInitialNodeCount = 1
)

type ClusterStatus string

const (
	ClusterStatusUnknown   ClusterStatus = "unknown"
	ClusterStatusUncreated ClusterStatus = "uncreated"
	ClusterStatusCreating  ClusterStatus = "creating"
	ClusterStatusReady     ClusterStatus = "ready"
	ClusterStatusFailed    ClusterStatus = "failed"
)

type ClusterSpec struct {
	Name        string `json:"name"`
	Description string `json:"Description"`
	Nodes       int32  `json:"nodes"`
	MachineType string `json:"machine_type"`
	Region      string `json:"region"`
	ProjectID   string `json:"project_id"`
}

// Cluster manages the desired state of a wanted cluster in GKE
// Prototype flow and rules:
// - try and load our stored state about the cluster
// - check if the cluster exists in GKE
// - if it doesn't exist, create it
// - block and wait for status change
// - if the new status is a failure,
//     then set Cluster.Status,
//     update our state db,
//     and return an error back to caller
// - if the new status is success,
//     then set Cluster.Status,
//     update our state db,
//     and return nil back to caller
type Cluster struct {
	Status ClusterStatus
	URL    *url.URL
	Spec   ClusterSpec

	client         *container.ClusterManagerClient
	resourceID     string
	stateManager   state.Manager
	gkeCluster     *containerpb.Cluster
	kubeconfig     *clientcmdapi.Config
	kubeconfigPath string
	logger         logging.Logger
}

func (c *Cluster) LookupRemote(ctx context.Context) errors.Error {
	req := &containerpb.GetClusterRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s",
			c.Spec.ProjectID, c.Spec.Region, c.Spec.Name),
	}

	resp, err := c.client.GetCluster(ctx, req)
	if err != nil {
		if status, ok := status.FromError(err); ok {
			if status.Code() == codes.NotFound {
				c.Status = ClusterStatusUncreated

				return nil
			}
		}

		return errors.NewGcpClusterReadError().WithCause(err)
	}

	c.Status = ClusterStatusUnknown

	if resp.Status == containerpb.Cluster_RUNNING {
		c.Status = ClusterStatusReady
	}

	c.gkeCluster = resp

	config, err := c.config(resp.Name, resp.Endpoint, resp.MasterAuth)
	if err != nil {
		return errors.NewGcpClusterReadError().WithCause(err)
	}

	c.kubeconfig = config

	tmpdir, err := ioutil.TempDir("", "nebula-gke-")
	if err != nil {
		return errors.NewGcpClusterReadError().WithCause(err)
	}

	c.kubeconfigPath = filepath.Join(tmpdir, "kubeconfig")

	if err := clientcmd.WriteToFile(*config, c.kubeconfigPath); err != nil {
		return errors.NewGcpClusterReadError().WithCause(err)
	}

	return nil
}

func (c *Cluster) KubeconfigPath() string {
	return c.kubeconfigPath
}

func (c *Cluster) config(name string, endpoint string, auth *containerpb.MasterAuth) (*clientcmdapi.Config, error) {
	config := clientcmdapi.NewConfig()

	cluster := clientcmdapi.NewCluster()

	ca, err := base64.StdEncoding.DecodeString(auth.ClusterCaCertificate)
	if err != nil {
		return nil, err
	}

	cluster.CertificateAuthorityData = ca

	cluster.Server = fmt.Sprintf("https://%s", endpoint)

	authInfo := clientcmdapi.NewAuthInfo()
	authInfo.Username = auth.Username
	authInfo.Password = auth.Password

	clientCert, err := base64.StdEncoding.DecodeString(auth.ClientCertificate)
	if err != nil {
		return nil, err
	}

	authInfo.ClientCertificateData = clientCert

	clientKey, err := base64.StdEncoding.DecodeString(auth.ClientKey)
	if err != nil {
		return nil, err
	}

	authInfo.ClientKeyData = clientKey

	context := clientcmdapi.NewContext()
	context.AuthInfo = name
	context.Cluster = name
	context.Namespace = "default"

	config.Clusters[name] = cluster
	config.AuthInfos[name] = authInfo
	config.Contexts[name] = context
	config.CurrentContext = name

	return config, nil
}

func (c *Cluster) SaveState(ctx context.Context) errors.Error {
	value, err := c.encode()
	if err != nil {
		return err
	}

	if err := c.stateManager.Save(&state.Resource{Name: c.resourceID, Value: json.RawMessage(value)}); err != nil {
		return err
	}

	return nil
}

func (c *Cluster) Sync(ctx context.Context) errors.Error {
	if c.Status == ClusterStatusUncreated {
		return c.create(ctx)
	}

	return nil
}

func (c *Cluster) create(ctx context.Context) errors.Error {
	req := &containerpb.CreateClusterRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", c.Spec.ProjectID, c.Spec.Region),
		Cluster: &containerpb.Cluster{
			Name:             c.Spec.Name,
			Description:      c.Spec.Description,
			Location:         c.Spec.Region,
			InitialNodeCount: c.Spec.Nodes,
			NodeConfig: &containerpb.NodeConfig{
				MachineType: c.Spec.MachineType,
			},
		},
	}

	// TODO handle response
	_, err := c.client.CreateCluster(ctx, req)
	if err != nil {
		return errors.NewWorkflowUnknownRuntimeError().WithCause(err).Bug()
	}

	if err := c.LookupRemote(ctx); err != nil {
		return err
	}

	hm := helm.NewHelmManager(c.KubeconfigPath(), c.logger)

	return hm.InitTiller(ctx)
}

func (c *Cluster) encode() ([]byte, errors.Error) {
	b, err := json.Marshal(c.Spec)
	if err != nil {
		return nil, errors.NewGcpClusterEncodingError().WithCause(err)
	}

	return b, nil
}

func NewCluster(rid string, sm state.Manager, spec ClusterSpec, logger logging.Logger) (*Cluster, errors.Error) {
	manager, err := container.NewClusterManagerClient(context.Background())
	if err != nil {
		return nil, errors.NewGcpClientCreateError().WithCause(err)
	}

	if spec.MachineType == "" {
		spec.MachineType = defaultMachineType
	}

	return &Cluster{
		Status:       ClusterStatusUnknown,
		Spec:         spec,
		client:       manager,
		resourceID:   rid,
		stateManager: sm,
		logger:       logger,
	}, nil
}

func NewClusterFromResourceID(rid string, sm state.Manager, logger logging.Logger) (*Cluster, errors.Error) {
	manager, merr := container.NewClusterManagerClient(context.Background())
	if merr != nil {
		return nil, errors.NewGcpClientCreateError().WithCause(merr)
	}

	cluster := Cluster{
		client:       manager,
		resourceID:   rid,
		stateManager: sm,
		logger:       logger,
	}

	r, err := sm.Load(rid)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(r.Value, &cluster.Spec); err != nil {
		return nil, errors.NewGcpClientCreateError().WithCause(err)
	}

	return &cluster, nil
}
