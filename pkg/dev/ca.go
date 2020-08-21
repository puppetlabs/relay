package dev

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	"github.com/puppetlabs/relay/pkg/cluster"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

const (
	jwtSigningKeysSecretName = "jwt-signing-keys"
)

type rsaPair struct {
	publicKey  []byte
	privateKey []byte
}

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

func (m *caManager) createSigningKeys(ctx context.Context) (*corev1.Secret, error) {
	pair, err := m.generateJWTSigningKeys()
	if err != nil {
		return nil, err
	}

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jwtSigningKeysSecretName,
			Namespace: "relay-system",
		},
		Data: map[string][]byte{
			"private-key.pem": pair.privateKey,
			"public-key.pem":  pair.publicKey,
		},
	}

	if err := m.cl.APIClient.Create(ctx, sec); err != nil {
		return nil, err
	}

	return sec, nil
}

func (m *caManager) generateJWTSigningKeys() (*rsaPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}

	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})

	return &rsaPair{privateKey: privPEM, publicKey: pubPEM}, nil
}

func newCAManager(cl *cluster.Client) *caManager {
	return &caManager{cl: cl}
}
