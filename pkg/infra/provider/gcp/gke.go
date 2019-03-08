package gcp

import (
	"context"
	"fmt"

	container "cloud.google.com/go/container/apiv1"
	"github.com/puppetlabs/nebula/pkg/errors"
	containerpb "google.golang.org/genproto/googleapis/container/v1"

	// "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
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
	Region      string
	ProjectID   string
}

type Cluster struct {
	Status ClusterStatus
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
			if status.GetCode() == code.Code_value["NOT_FOUND"] {
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
	return c.create(ctx)
}

func (c *Cluster) create(ctx context.Context) errors.Error {
	req := &containerpb.CreateClusterRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", c.Spec.ProjectID, c.Spec.Region),
		Cluster: &containerpb.Cluster{
			Name:             c.Spec.Name,
			InitialNodeCount: c.Spec.Nodes,
			Description:      c.Spec.Description,
			Location:         c.Spec.Region,
		},
	}

	resp, err := c.client.CreateCluster(ctx, req)
	if err != nil {
		return errors.NewWorkflowUnknownRuntimeError().WithCause(err).Bug()
	}

	return nil
}

func NewCluster(spec ClusterSpec) (*Cluster, errors.Error) {
	manager, err := container.NewClusterManagerClient(context.Background())
	if err != nil {
		return nil, errors.NewGcpClientCreateError().WithCause(err)
	}

	return &Cluster{
		Status: ClusterStatusUnknown,
		Spec:   spec,
		client: manager,
	}, nil
}
