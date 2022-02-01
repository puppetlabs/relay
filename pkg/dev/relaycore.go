package dev

import (
	"context"
	"errors"
	"fmt"

	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	installerv1alpha1 "github.com/puppetlabs/relay-core/pkg/apis/install.relay.sh/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	relayCoreName                  = "relay-core-v1"
	relayOperatorStorageVolumeSize = "1Gi"
)

type relayCoreObjects struct {
	pvc                corev1.PersistentVolumeClaim
	relayCore          installerv1alpha1.RelayCore
	serviceAccount     corev1.ServiceAccount
	clusterRoleBinding rbacv1.ClusterRoleBinding
}

func newRelayCoreObjects() *relayCoreObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      relayCoreName,
		Namespace: systemNamespace,
	}

	operatorObjectMeta := objectMeta
	operatorObjectMeta.Name = fmt.Sprintf("%s-operator", objectMeta.Name)

	operatorAdminObjectMeta := objectMeta
	operatorAdminObjectMeta.Name = fmt.Sprintf("%s-operator-admin", objectMeta.Name)

	return &relayCoreObjects{
		pvc:                corev1.PersistentVolumeClaim{ObjectMeta: operatorObjectMeta},
		relayCore:          installerv1alpha1.RelayCore{ObjectMeta: objectMeta},
		serviceAccount:     corev1.ServiceAccount{ObjectMeta: operatorObjectMeta},
		clusterRoleBinding: rbacv1.ClusterRoleBinding{ObjectMeta: operatorAdminObjectMeta},
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

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.pvc, func() error {
		m.operatorStoragePVC(&m.objects.pvc)

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

	if _, err := ctrl.CreateOrUpdate(ctx, cl, &m.objects.clusterRoleBinding, func() error {
		m.rbacDefinition(&m.objects.clusterRoleBinding)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *relayCoreManager) operatorStoragePVC(pvc *corev1.PersistentVolumeClaim) {
	pvc.Spec.AccessModes = []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce}

	pvc.Spec.Resources = corev1.ResourceRequirements{
		Requests: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceStorage: resource.MustParse(relayOperatorStorageVolumeSize),
		},
	}
}

func (m *relayCoreManager) relayCore(rc *installerv1alpha1.RelayCore) {
	if m.logServiceOpts.Enabled {
		rc.Spec.LogService = installerv1alpha1.LogServiceConfig{
			Image:           m.installerOpts.LogServiceImage,
			ImagePullPolicy: corev1.PullAlways,
			CredentialsSecretKeyRef: corev1.SecretKeySelector{
				Key: m.logServiceOpts.CredentialsKey,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: m.logServiceOpts.CredentialsSecretName,
				},
			},
			Project: m.logServiceOpts.Project,
			Dataset: m.logServiceOpts.Dataset,
			Table:   m.logServiceOpts.Table,
		}
	}

	rc.Spec.Operator = &installerv1alpha1.OperatorConfig{
		Image:             m.installerOpts.OperatorImage,
		ImagePullPolicy:   corev1.PullAlways,
		Standalone:        true,
		LogStoragePVCName: &m.objects.pvc.Name,
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

	rc.Spec.MetadataAPI = &installerv1alpha1.MetadataAPIConfig{
		Image:           m.installerOpts.MetadataAPIImage,
		ImagePullPolicy: corev1.PullAlways,
		VaultAuthRole:   "tenant",
		VaultAuthPath:   "auth/jwt-tenants",
	}

	rc.Spec.VaultDeployment = &installerv1alpha1.VaultDeployment{
		Image:           "gcr.io/nebula-tasks/nebula-vault:1.7.3-oauthapp-3.0.0-beta.3-2.2.0-1.10.0",
		ImagePullPolicy: corev1.PullAlways,
		Configuration: `
  disable_mlock = true
  ui = true
  plugin_directory = "/nebula/vault/plugins"
  log_level = "Debug"
  listener "tcp" {
    tls_disable = 1
	address = "0.0.0.0:8200"
    // address = "[::]:8200"
    // cluster_address = "[::]:8201"
  }
  storage "file" {
    path = "/vault/data"
  }`,
	}

	rc.Spec.Vault = &installerv1alpha1.VaultConfig{
		VaultInitializationImage:           m.installerOpts.OperatorVaultInitImage,
		VaultInitializationImagePullPolicy: corev1.PullAlways,

		// FIXME Change this to be more flexible/specific
		AuthDelegatorServiceAccountName: vaultIdentifier,

		LogServicePath: "pls",
		TenantPath:     "customers",
		TransitKey:     "metadata-api",
		TransitPath:    "transit-tenants",
	}
}

func (m *relayCoreManager) rbacDefinition(crb *rbacv1.ClusterRoleBinding) {
	crb.RoleRef =
		rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		}
	crb.Subjects = []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      m.objects.serviceAccount.Name,
			Namespace: m.objects.serviceAccount.Namespace,
		},
	}

}

func (m *relayCoreManager) wait(ctx context.Context) error {
	err := retry.Wait(ctx, func(ctx context.Context) (bool, error) {
		key := client.ObjectKeyFromObject(&m.objects.relayCore)

		if err := m.cl.APIClient.Get(ctx, key, &m.objects.relayCore); err != nil {
			return retry.Repeat(err)
		}

		if m.objects.relayCore.Status.Status != installerv1alpha1.StatusCreated {
			return retry.Repeat(errors.New("waiting for relaycore to be created"))
		}

		return retry.Done(nil)
	})
	if err != nil {
		return err
	}

	return nil
}

func newRelayCoreManager(cl *Client, installerOpts InstallerOptions, logServiceOpts LogServiceOptions) *relayCoreManager {
	return &relayCoreManager{
		cl:             cl,
		objects:        newRelayCoreObjects(),
		installerOpts:  installerOpts,
		logServiceOpts: logServiceOpts,
	}
}
