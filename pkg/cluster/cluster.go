package cluster

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const (
	ClusterName = "relay-workflows"
	NetworkName = "relay-workflows-net"
	WorkerCount = 2
)

var agentArgs = []string{
	"--node-label=nebula.puppet.com/scheduling.customer-ready=true",
}

type ClientOptions struct {
	Scheme *runtime.Scheme
}

// Manager provides methods to manage the lifecycle of a cluster.
type Manager interface {
	Exists(ctx context.Context) (bool, error)
	Create(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Delete(ctx context.Context) error
	GetKubeconfig(ctx context.Context) (*clientcmdapi.Config, error)
	WriteKubeconfig(ctx context.Context, path string) error
	GetClient(ctx context.Context, opts ClientOptions) (*Client, error)
}

// NewManager returns a new selected Manager. Since k3d is
// the only one supported right now, we just return a manager
// that one. Using this function ensures the delegate manager
// always satisfies the Manager interface.
func NewManager() Manager {
	return NewK3dClusterManager()
}
