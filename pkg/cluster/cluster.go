package cluster

import (
	"context"

	"github.com/puppetlabs/leg/workdir"
	"k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	ClusterName               = "relay-workflows"
	NetworkName               = "relay-workflows-net"
	HostStorageName           = "local-storage"
	ImagePassthroughCacheAddr = "http://localhost:5001"

	DefaultRegistryName         = "docker-registry.docker-registry.svc.cluster.local"
	DefaultRegistryPort         = 5000
	DefaultLoadBalancerHostPort = 8080
	DefaultLoadBalancerNodePort = 80
	DefaultWorkerCount          = 0
)

type ClientOptions struct {
	Scheme *runtime.Scheme
}

type Config struct {
	WorkDir *workdir.WorkDir
}

// CreateOptions are the configurable options for cluster creation
type CreateOptions struct {
	// LoadBalancerHostPort is the port on the host to bind to when mapping
	// between the host machine and the service load balancer node.
	LoadBalancerHostPort int
	// Number of worker nodes on the cluster
	WorkerCount int
}

// InitializeOptions are the configurable options for cluster initialization
type InitializeOptions struct{}

// Manager provides methods to manage the lifecycle of a cluster.
type Manager interface {
	Exists(ctx context.Context) (bool, error)
	Create(ctx context.Context, opts CreateOptions) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Delete(ctx context.Context) error
	GetKubeconfig(ctx context.Context) (*clientcmdapi.Config, error)
	WriteKubeconfig(ctx context.Context) error
}

// NewManager returns a new selected Manager. Since k3d is
// the only one supported right now, we just return a manager
// that one. Using this function ensures the delegate manager
// always satisfies the Manager interface.
func NewManager(cfg Config) Manager {
	return NewK3dClusterManager(cfg)
}
