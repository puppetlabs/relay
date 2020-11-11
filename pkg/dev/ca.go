package dev

import (
	"github.com/puppetlabs/relay/pkg/cluster"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	jwtSigningKeysSecretName = "relay-core-v1-operator-signing-keys"
)

type caManager struct {
	cl *cluster.Client
}

func (m *caManager) admissionPatcher(caCertData []byte) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *admissionregistrationv1beta1.MutatingWebhookConfiguration:
			for i := range t.Webhooks {
				t.Webhooks[i].ClientConfig.CABundle = caCertData
			}
		}
	}
}

func newCAManager(cl *cluster.Client) *caManager {
	return &caManager{cl: cl}
}
