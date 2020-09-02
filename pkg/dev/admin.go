package dev

import (
	"context"
	"fmt"
	"path"

	"github.com/puppetlabs/relay/pkg/cluster"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	relayAdminServiceAccountName = "relay-admin-user"
	relayClusterConnectionID     = "_relay-dev-cluster"

	RelayClusterConnectionName = "relay-dev-cluster"
)

type adminManager struct {
	cl *cluster.Client
	vm *vaultManager
}

func (m *adminManager) createServiceAccount(ctx context.Context) error {
	sa := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      relayAdminServiceAccountName,
			Namespace: "kube-system",
		},
	}

	crb := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: relayAdminServiceAccountName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      relayAdminServiceAccountName,
				Namespace: "kube-system",
			},
		},
	}

	if err := m.cl.APIClient.Create(ctx, sa); err != nil {
		return fmt.Errorf("failed to create admin service account: %w", err)
	}

	if err := m.cl.APIClient.Create(ctx, crb); err != nil {
		return fmt.Errorf("failed to create admin user cluster role binding: %w", err)
	}

	return nil
}

func (m *adminManager) addConnectionForWorkflow(ctx context.Context, name string) error {
	saKey := client.ObjectKey{Name: relayAdminServiceAccountName, Namespace: "kube-system"}
	sa := &corev1.ServiceAccount{}
	saSecret := &corev1.Secret{}

	if err := m.cl.APIClient.Get(ctx, saKey, sa); err != nil {
		return err
	}

	secretKey := client.ObjectKey{Name: sa.Secrets[0].Name, Namespace: "kube-system"}

	if err := m.cl.APIClient.Get(ctx, secretKey, saSecret); err != nil {
		return err
	}

	connectionsPath := path.Join("customers", "connections", name)
	pointerPath := path.Join(connectionsPath, "kubernetes", RelayClusterConnectionName)
	base := path.Join(connectionsPath, relayClusterConnectionID)
	connectionSecrets := map[string]string{
		pointerPath:                             relayClusterConnectionID,
		path.Join(base, "token"):                string(saSecret.Data["token"]),
		path.Join(base, "certificateAuthority"): string(saSecret.Data["ca.crt"]),
		path.Join(base, "server"):               "https://kubernetes.default.svc.cluster.local",
	}

	return m.vm.writeSecrets(ctx, connectionSecrets)
}

func newAdminManager(cl *cluster.Client, vm *vaultManager) *adminManager {
	return &adminManager{
		cl: cl,
		vm: vm,
	}
}
