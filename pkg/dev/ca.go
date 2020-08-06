package dev

import (
	"context"

	certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	certmanagermetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/puppetlabs/relay/pkg/cluster"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type caPair struct {
	certificate []byte
	privateKey  []byte
}

type caManager struct {
	cl *cluster.Client
}

func (m *caManager) secretPatcher(target string, pair *caPair) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *corev1.Secret:
			if t.GetName() == target {
				t.Data = map[string][]byte{
					"tls.crt": pair.certificate,
					"tls.key": pair.privateKey,
				}
			}
		}
	}
}

func (m *caManager) certificatePatcher(issuerRef string) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *certmanagerv1beta1.Certificate:
			t.Spec.IssuerRef = certmanagermetav1.ObjectReference{
				Name: issuerRef,
				Kind: "ClusterIssuer",
			}
		}
	}
}

func (m *caManager) admissionPatcher(caSecret client.ObjectKey) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *admissionregistrationv1beta1.MutatingWebhookConfiguration:
			tlsSecret := &corev1.Secret{}

			_ = m.cl.APIClient.Get(context.TODO(), caSecret, tlsSecret)

			for i := range t.Webhooks {
				t.Webhooks[i].ClientConfig.CABundle = tlsSecret.Data["ca.crt"]
			}
		}
	}
}

func newCAManager(cl *cluster.Client) *caManager {
	return &caManager{cl: cl}
}
