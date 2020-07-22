package cluster

import (
	"context"

	"github.com/puppetlabs/relay/pkg/config"
)

const (
	DefaultClusterName = "relay-workflows"
	DefaultNetworkName = "relay-workflows-net"
	DefaultWorkerCount = 3
	DefaultK3sVersion  = "latest"
)

type Manager interface {
	Exists(ctx context.Context) (bool, error)
	Create(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Delete(ctx context.Context) error
}

// NewManager returns a new selected Manager. Since k3d is
// the only one supported right now, we just return a manager
// that one. Using this function ensures the delegate manager
// always satisfies the Manager interface.
func NewManager(cfg *config.Config) Manager {
	return NewK3dClusterManager(cfg)
}
