package dev

import (
	"context"
	"fmt"

	installerv1alpha1 "github.com/puppetlabs/relay-core/pkg/apis/install.relay.sh/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	RelayInstallerImage                            = "us-docker.pkg.dev/puppet-relay-contrib-oss/relay-core/relay-installer"
	RelayMetadataAPIImage                          = "us-docker.pkg.dev/puppet-relay-contrib-oss/relay-core/relay-metadata-api"
	RelayOperatorImage                             = "us-docker.pkg.dev/puppet-relay-contrib-oss/relay-core/relay-operator"
	RelayOperatorVaultInitImage                    = "us-docker.pkg.dev/puppet-relay-contrib-oss/relay-core/relay-operator-vault-init"
	RelayOperatorWebhookCertificateControllerImage = "us-docker.pkg.dev/puppet-relay-contrib-oss/relay-core/relay-operator-webhook-certificate-controller"

	RelayLogServiceImage = "relaysh/relay-pls:latest"
)

const (
	DefaultVaultConfiguration = `
disable_mlock = true
ui = true
log_level = "Debug"
listener "tcp" {
	tls_disable = 1
	address = "0.0.0.0:8200"
}
plugin_directory = "/relay/vault/plugins"
storage "file" {
	path = "/vault/data"
}`
	DefaultVaultConfigurationFile = "vault.hcl"
	DefaultVaultServerImage       = "relaysh/relay-vault:latest"
	DefaultVaultSidecarImage      = "vault:latest"
)

const (
	relayCoreName = "relay-core-v1"
)

type relayCoreObjects struct {
	configMap      corev1.ConfigMap
	relayCore      installerv1alpha1.RelayCore
	serviceAccount corev1.ServiceAccount
}

func newRelayCoreObjects() *relayCoreObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      relayCoreName,
		Namespace: systemNamespace,
	}

	operatorObjectMeta := objectMeta
	operatorObjectMeta.Name = fmt.Sprintf("%s-operator", objectMeta.Name)

	return &relayCoreObjects{
		configMap:      corev1.ConfigMap{ObjectMeta: operatorObjectMeta},
		relayCore:      installerv1alpha1.RelayCore{ObjectMeta: objectMeta},
		serviceAccount: corev1.ServiceAccount{ObjectMeta: operatorObjectMeta},
	}
}

type relayCoreManager struct {
	cl             *Client
	objects        *relayCoreObjects
	installerOpts  InstallerOptions
	logServiceOpts LogServiceOptions
}

func (m *relayCoreManager) reconcile(ctx context.Context) error {
	cl := m.cl.APIClient

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.configMap, func() error {
		m.operatorConfigMap(&m.objects.configMap)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.relayCore, func() error {
		m.relayCore(&m.objects.relayCore)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *relayCoreManager) operatorConfigMap(configMap *corev1.ConfigMap) {
	configMap.Data = map[string]string{
		DefaultVaultConfigurationFile: DefaultVaultConfiguration,
	}
}

func (m *relayCoreManager) relayCore(rc *installerv1alpha1.RelayCore) {
	rc.Spec.LogService = &installerv1alpha1.LogServiceConfig{
		Image:           m.installerOpts.LogServiceImage,
		ImagePullPolicy: corev1.PullAlways,
	}

	if m.logServiceOpts.CredentialsSecretName != "" && m.logServiceOpts.CredentialsKey != "" {
		rc.Spec.LogService.CredentialsSecretKeyRef = &corev1.SecretKeySelector{
			Key: m.logServiceOpts.CredentialsKey,
			LocalObjectReference: corev1.LocalObjectReference{
				Name: m.logServiceOpts.CredentialsSecretName,
			},
		}
		rc.Spec.LogService.Project = m.logServiceOpts.Project
		rc.Spec.LogService.Dataset = m.logServiceOpts.Dataset
		rc.Spec.LogService.Table = m.logServiceOpts.Table
	}

	tn := tenantNamespace
	rc.Spec.Operator = installerv1alpha1.OperatorConfig{
		Image:           m.installerOpts.OperatorImage,
		ImagePullPolicy: corev1.PullAlways,
		TenantNamespace: &tn,
		Standalone:      true,
		AdmissionWebhookServer: &installerv1alpha1.AdmissionWebhookServerConfig{
			CertificateControllerImage:           m.installerOpts.OperatorWebhookCertificateControllerImage,
			CertificateControllerImagePullPolicy: corev1.PullAlways,
			Domain:                               "admission.controller.relay.sh",
			NamespaceSelector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"controller.relay.sh/tenant-workload": "true",
				},
			},
		},
	}

	rc.Spec.MetadataAPI = installerv1alpha1.MetadataAPIConfig{
		Image:           m.installerOpts.MetadataAPIImage,
		ImagePullPolicy: corev1.PullAlways,
	}

	rc.Spec.Vault = installerv1alpha1.VaultConfig{
		Engine: installerv1alpha1.VaultEngineConfig{
			VaultInitializationImage:           m.installerOpts.OperatorVaultInitImage,
			VaultInitializationImagePullPolicy: corev1.PullAlways,

			// FIXME Change this to be more flexible/specific
			AuthDelegatorServiceAccountName: vaultIdentifier,
		},
		Server: installerv1alpha1.VaultServerConfig{
			BuiltIn: &installerv1alpha1.VaultServerBuiltInConfig{
				Image:           m.installerOpts.VaultServerImage,
				ImagePullPolicy: corev1.PullAlways,
				Resources: corev1.ResourceRequirements{
					Limits: map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceCPU:    resource.MustParse("50m"),
						corev1.ResourceMemory: resource.MustParse("64Mi"),
					},
					Requests: map[corev1.ResourceName]resource.Quantity{
						corev1.ResourceCPU:    resource.MustParse("25m"),
						corev1.ResourceMemory: resource.MustParse("64Mi"),
					},
				},
				ConfigMapRef: corev1.LocalObjectReference{
					Name: m.objects.configMap.Name,
				},
			},
		},
		Sidecar: installerv1alpha1.VaultSidecarConfig{
			Image:           m.installerOpts.VaultSidecarImage,
			ImagePullPolicy: corev1.PullAlways,
			Resources: corev1.ResourceRequirements{
				Limits: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("50m"),
					corev1.ResourceMemory: resource.MustParse("64Mi"),
				},
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("25m"),
					corev1.ResourceMemory: resource.MustParse("64Mi"),
				},
			},
		},
	}
}

func newRelayCoreManager(cl *Client, installerOpts InstallerOptions, logServiceOpts LogServiceOptions) *relayCoreManager {
	return &relayCoreManager{
		cl:             cl,
		objects:        newRelayCoreObjects(),
		installerOpts:  installerOpts,
		logServiceOpts: logServiceOpts,
	}
}
