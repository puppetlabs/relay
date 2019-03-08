package gcp

import (
	"context"
	"fmt"
	"net/url"

	container "cloud.google.com/go/container/apiv1"
	"github.com/puppetlabs/nebula/pkg/errors"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	Name        string
	Description string
	Nodes       int32
	MachineType string
	Region      string
	ProjectID   string
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

	client *container.ClusterManagerClient
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

	return c.LookupRemote(ctx)
}

func NewCluster(spec ClusterSpec) (*Cluster, errors.Error) {
	manager, err := container.NewClusterManagerClient(context.Background())
	if err != nil {
		return nil, errors.NewGcpClientCreateError().WithCause(err)
	}

	if spec.MachineType == "" {
		spec.MachineType = defaultMachineType
	}

	return &Cluster{
		Status: ClusterStatusUnknown,
		Spec:   spec,
		client: manager,
	}, nil
}
