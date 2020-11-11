package dev

import (
	"context"
	"path"

	"github.com/puppetlabs/relay/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	relayAdminServiceAccountName = "relay-admin-user"
	relayClusterConnectionID     = "_relay-dev-cluster"

	RelayClusterConnectionName = "relay-dev-cluster"
)

type adminObjects struct {
	serviceAccount     corev1.ServiceAccount
	secret             corev1.Secret
	clusterRoleBinding rbacv1.ClusterRoleBinding
}

func newAdminObjects() *adminObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      relayAdminServiceAccountName,
		Namespace: systemNamespace,
	}

	return &adminObjects{
		serviceAccount: corev1.ServiceAccount{ObjectMeta: objectMeta},
		secret:         corev1.Secret{ObjectMeta: objectMeta},
		clusterRoleBinding: rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: objectMeta.Name,
			},
		},
	}
}

type adminManager struct {
	cl      *cluster.Client
	objects *adminObjects
	vm      *vaultManager
}

func (m *adminManager) reconcile(ctx context.Context) error {
	client := m.cl.APIClient

	if _, err := ctrl.CreateOrUpdate(ctx, client, &m.objects.serviceAccount, func() error {
		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, client, &m.objects.secret, func() error {
		m.secret(&m.objects.secret)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, client, &m.objects.clusterRoleBinding, func() error {
		m.clusterRoleBinding(&m.objects.clusterRoleBinding)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *adminManager) secret(secret *corev1.Secret) {
	if secret.Annotations == nil {
		secret.Annotations = make(map[string]string)
	}

	secret.Annotations["kubernetes.io/service-account.name"] = m.objects.serviceAccount.Name
	secret.Type = corev1.SecretTypeServiceAccountToken
}

func (m *adminManager) clusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) {
	clusterRoleBinding.RoleRef = rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "cluster-admin",
	}

	clusterRoleBinding.Subjects = []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      relayAdminServiceAccountName,
			Namespace: systemNamespace,
		},
	}
}

func (m *adminManager) addConnectionForWorkflow(ctx context.Context, name string) error {
	secretKey, err := client.ObjectKeyFromObject(&m.objects.secret)
	if err != nil {
		return err
	}

	if err := m.cl.APIClient.Get(ctx, secretKey, &m.objects.secret); err != nil {
		return err
	}

	data := m.objects.secret.Data

	connectionsPath := path.Join("customers", "connections", name)
	pointerPath := path.Join(connectionsPath, "kubernetes", RelayClusterConnectionName)
	base := path.Join(connectionsPath, relayClusterConnectionID)

	connectionSecrets := map[string]string{
		pointerPath:                             relayClusterConnectionID,
		path.Join(base, "token"):                string(data["token"]),
		path.Join(base, "certificateAuthority"): string(data["ca.crt"]),
		path.Join(base, "server"):               "https://kubernetes.default.svc.cluster.local",
	}

	return m.vm.writeSecrets(ctx, connectionSecrets)
}

func newAdminManager(cl *cluster.Client, vm *vaultManager) *adminManager {
	return &adminManager{
		cl:      cl,
		objects: newAdminObjects(),
		vm:      vm,
	}
}
