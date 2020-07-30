package dev

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	systemNamespace    = "relay-system"
	workflowsNamespace = "relay-workflows"
)

type namespaceManager struct {
	cl client.Client
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

	if err := m.cl.Create(ctx, sn); err != nil {
		return err
	}

	if err := m.cl.Create(ctx, wn); err != nil {
		return err
	}

	return nil
}

func (m *namespaceManager) objectNamespacePatcher(id string) objectPatcherFunc {
	return func(obj runtime.Object) {
		_ := obj.GetObjectKind().GroupVersionKind()

		a, err := meta.Accessor(obj)
		if err != nil {
			return
		}

		var ns string

		switch id {
		case "system":
			ns = systemNamespace
		case "workflows":
			ns = workflowsNamespace
		default:
			return
		}

		if a.GetNamespace() == "" {
			a.SetNamespace(ns)
		}
	}
}

func newNamespaceManager(cl client.Client) *namespaceManager {
	return &namespaceManager{cl: cl}
}
