package dev

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	systemNamespace = "relay-system"
	tenantNamespace = "relay-tenants"

	knativeServingNamespace  = "knative-serving"
	kourierSystemNamespace   = "kourier-system"
	tektonPipelinesNamespace = "tekton-pipelines"
)

type namespaceObjects struct {
	systemNamespace         corev1.Namespace
	knativeServingNamespace corev1.Namespace
	kourierSystemNamespace  corev1.Namespace
}

func newNamespaceObjects() *namespaceObjects {
	return &namespaceObjects{
		systemNamespace:         corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: systemNamespace}},
		knativeServingNamespace: corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: knativeServingNamespace}},
		kourierSystemNamespace:  corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: kourierSystemNamespace}},
	}
}

type namespaceManager struct {
	cl      *Client
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

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.knativeServingNamespace, func() error {
		m.knativeServingNamespace(&m.objects.knativeServingNamespace)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.knativeServingNamespace, func() error {
		m.kourierSystemNamespace(&m.objects.kourierSystemNamespace)

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

func (m *namespaceManager) knativeServingNamespace(ns *corev1.Namespace) {
	ns.Labels = map[string]string{
		"nebula.puppet.com/network-policy.webhook-gateway": "true",
	}
}

func (m *namespaceManager) kourierSystemNamespace(ns *corev1.Namespace) {
	ns.Labels = map[string]string{
		"nebula.puppet.com/network-policy.webhook-gateway": "true",
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

func newNamespaceManager(cl *Client) *namespaceManager {
	return &namespaceManager{
		cl:      cl,
		objects: newNamespaceObjects(),
	}
}
