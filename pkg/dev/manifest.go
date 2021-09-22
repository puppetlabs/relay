package dev

import (
	"context"

	"github.com/puppetlabs/leg/k8sutil/pkg/manifest"
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev/manifests"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ManifestManager struct {
	cl *cluster.Client
}

func (m *ManifestManager) ProcessManifests(ctx context.Context, path string, patchers ...manifest.PatcherFunc) error {
	files := manifests.MustAssetListDir(path)

	for _, file := range files {
		r := manifests.MustAsset(file)

		objs, err := manifest.Parse(DefaultScheme, r, patchers...)
		if err != nil {
			return nil
		}

		for _, obj := range objs {
			if err := m.cl.APIClient.Patch(ctx, obj, client.Apply, client.ForceOwnership, client.FieldOwner("relay")); err != nil {
				return err
			}
		}
	}

	return nil
}

func NewManifestManager(cl *cluster.Client) *ManifestManager {
	return &ManifestManager{
		cl: cl,
	}
}
