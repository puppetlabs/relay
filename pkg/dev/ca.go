package dev

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"time"

	certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	certmanagermetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type caPair struct {
	certificate []byte
	privateKey  []byte
}

type caManager struct{}

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

func (m *caManager) admissionPatcher(pair *caPair) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *admissionregistrationv1beta1.MutatingWebhookConfiguration:
			for i := range t.Webhooks {
				t.Webhooks[i].ClientConfig.CABundle = []byte(base64.StdEncoding.EncodeToString(pair.certificate))
			}
		}
	}
}

func (m *caManager) generateCA() (*caPair, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization: []string{""},
			Country:      []string{"US"},
			Province:     []string{},
			Locality:     []string{"Portland"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return nil, err
	}

	caPEM := &bytes.Buffer{}
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := &bytes.Buffer{}
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	return &caPair{
		certificate: caPEM.Bytes(),
		privateKey:  caPrivKeyPEM.Bytes(),
	}, nil
}

func newCAManager() *caManager {
	return &caManager{}
}
