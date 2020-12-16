package dev

import (
	"context"

	"github.com/puppetlabs/relay/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	systemNamespace          = "relay-system"
	tenantNamespace          = "relay-tenants"
	registryNamespace        = "docker-registry"
	ambassadorNamespace      = "ambassador-webhook"
	knativeServingNamespace  = "knative-serving"
	tektonPipelinesNamespace = "tekton-pipelines"
)

type namespaceObjects struct {
	systemNamespace         corev1.Namespace
	tenantNamespace         corev1.Namespace
	registryNamespace       corev1.Namespace
	ambassadorNamespace     corev1.Namespace
	knativeServingNamespace corev1.Namespace
}

func newNamespaceObjects() *namespaceObjects {
	return &namespaceObjects{
		systemNamespace:         corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: systemNamespace}},
		tenantNamespace:         corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: tenantNamespace}},
		registryNamespace:       corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: registryNamespace}},
		ambassadorNamespace:     corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ambassadorNamespace}},
		knativeServingNamespace: corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: knativeServingNamespace}},
	}
}

type namespaceManager struct {
	cl      *cluster.Client
	objects *namespaceObjects
}

func (m *namespaceManager) reconcile(ctx context.Context) error {
	cl := m.cl.APIClient

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.systemNamespace, func() error {
		m.systemNamespace(&m.objects.systemNamespace)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.tenantNamespace, func() error {
		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.registryNamespace, func() error {
		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.ambassadorNamespace, func() error {
		m.ambassadorNamespace(&m.objects.ambassadorNamespace)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.knativeServingNamespace, func() error {
		m.knativeServingNamespace(&m.objects.knativeServingNamespace)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *namespaceManager) systemNamespace(ns *corev1.Namespace) {
	ns.Labels = map[string]string{
		"nebula.puppet.com/network-policy.tasks": "true",
	}
}

func (m *namespaceManager) ambassadorNamespace(ns *corev1.Namespace) {
	ns.Labels = map[string]string{
		"nebula.puppet.com/network-policy.webhook-gateway": "true",
	}
}

func (m *namespaceManager) knativeServingNamespace(ns *corev1.Namespace) {
	ns.Labels = map[string]string{
		"nebula.puppet.com/network-policy.webhook-gateway": "true",
	}
}

func (m *namespaceManager) objectNamespacePatcher(name string) objectPatcherFunc {
	return func(obj runtime.Object) {
		var gvk schema.GroupVersionKind

		gvks, _, err := DefaultScheme.ObjectKinds(obj)
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

		if a.GetNamespace() == "" {
			if name == "" {
				a.SetNamespace("default")
			} else {
				a.SetNamespace(name)
			}
		}
	}
}

func (m *namespaceManager) delete(ctx context.Context, ns string) error {
	namespace := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: ns,
		},
	}

	return m.cl.APIClient.Delete(ctx, namespace, client.PropagationPolicy(metav1.DeletePropagationBackground))
}

func newNamespaceManager(cl *cluster.Client) *namespaceManager {
	return &namespaceManager{
		cl:      cl,
		objects: newNamespaceObjects(),
	}
}
