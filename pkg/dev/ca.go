package dev

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"math/big"
	"time"

	"github.com/puppetlabs/relay/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	caSecretName = "relay-ca"
)

type caPair struct {
	certificate []byte
	privateKey  []byte
}

type caManager struct {
	cl *cluster.Client

	patchers []objectPatcherFunc
}

func (m *caManager) create(ctx context.Context) error {
	pair, err := m.generateCA()
	if err != nil {
		return err
	}

	caSecret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: caSecretName,
		},
		StringData: map[string]string{
			"tls.crt": base64.StdEncoding.EncodeToString(pair.certificate),
			"tls.key": base64.StdEncoding.EncodeToString(pair.privateKey),
		},
	}

	for _, patcher := range m.patchers {
		patcher(caSecret)
	}

	if err := m.cl.APIClient.Create(ctx, caSecret); err != nil {
		return err
	}

	return nil
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

func newCAManager(cl *cluster.Client, patchers ...objectPatcherFunc) *caManager {
	return &caManager{
		cl:       cl,
		patchers: patchers,
	}
}
