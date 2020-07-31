package dev

import (
	"context"

	"github.com/puppetlabs/relay/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	systemNamespace    = "relay-system"
	workflowsNamespace = "relay-workflows"
)

type namespaceManager struct {
	cl *cluster.Client
}

func (m *namespaceManager) create(ctx context.Context) error {
	sn := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: systemNamespace,
		},
	}

	wn := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: workflowsNamespace,
		},
	}

	if err := m.cl.APIClient.Create(ctx, sn); err != nil {
		return err
	}

	if err := m.cl.APIClient.Create(ctx, wn); err != nil {
		return err
	}

	return nil
}

func (m *namespaceManager) getByID(id string) string {
	switch id {
	case "system":
		return systemNamespace
	case "workflows":
		return workflowsNamespace
	}

	return "default"
}

func (m *namespaceManager) objectNamespacePatcher(id string) objectPatcherFunc {
	return func(obj runtime.Object) {
		var gvk schema.GroupVersionKind

		gvks, _, err := scheme.ObjectKinds(obj)
		if err != nil {
			return
		}

		if len(gvks) > 1 {
			return
		}

		gvk = gvks[0]

		mapping, err := m.cl.Mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return
		}

		a, err := meta.Accessor(obj)
		if err != nil {
			return
		}

		if mapping.Scope.Name() != meta.RESTScopeNameNamespace {
			return
		}

		ns := m.getByID(id)

		if a.GetNamespace() == "" {
			a.SetNamespace(ns)
		}
	}
}

func newNamespaceManager(cl *cluster.Client) *namespaceManager {
	return &namespaceManager{cl: cl}
}
