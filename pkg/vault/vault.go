package vault

import (
	"context"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/puppetlabs/relay/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type VaultSystemClient struct {
	kubeClient  client.Client
	vaultClient *vaultapi.Client
}

func (vsc *VaultSystemClient) InitializeVault(ctx context.Context, key types.NamespacedName, plugins []*vaultapi.RegisterPluginInput) error {
	credentials, err := vsc.getCredentials(ctx, key)
	if err != nil {
		return err
	}

	if credentials == nil {
		credentials, err = vsc.initialize(ctx, key)
		if err != nil {
			return err
		}
	}

	err = vsc.unseal(credentials)
	if err != nil {
		return err
	}

	vsc.vaultClient.SetToken(credentials.RootToken)

	err = vsc.registerPlugins(plugins)
	if err != nil {
		return err
	}

	err = vsc.enableKubernetesAuth(ctx, key)
	if err != nil {
		return err
	}

	err = vsc.configureKubernetesAuth(ctx, key)
	if err != nil {
		return err
	}

	return nil
}

func (vsc *VaultSystemClient) initialize(ctx context.Context, key types.NamespacedName) (*VaultKeys, error) {
	response, err := vsc.vaultClient.Sys().Init(
		&vaultapi.InitRequest{
			SecretShares:    1,
			SecretThreshold: 1,
		},
	)
	if err != nil {
		return nil, err
	}

	vaultKeys := &VaultKeys{
		RootToken:  response.RootToken,
		UnsealKeys: response.Keys,
	}

	err = vsc.createCredentials(ctx, key, vaultKeys)
	if err != nil {
		return nil, err
	}

	return vaultKeys, nil
}

func (vsc *VaultSystemClient) getCredentials(ctx context.Context, key types.NamespacedName) (*VaultKeys, error) {
	secret := &corev1.Secret{}
	err := vsc.kubeClient.Get(ctx, key, secret)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			return nil, nil
		}

		return nil, err
	}

	return &VaultKeys{
		RootToken:  string(secret.Data[VaultRootToken]),
		UnsealKeys: []string{string(secret.Data[VaultUnsealKey])},
	}, nil
}

func (vsc *VaultSystemClient) createCredentials(ctx context.Context, key types.NamespacedName, vaultKeys *VaultKeys) error {
	objectMeta := metav1.ObjectMeta{
		Name:      key.Name,
		Namespace: key.Namespace,
	}

	secret := &corev1.Secret{
		ObjectMeta: objectMeta,
		StringData: map[string]string{
			VaultRootToken: vaultKeys.RootToken,
			VaultUnsealKey: vaultKeys.UnsealKeys[0],
		},
	}

	if err := vsc.kubeClient.Create(ctx, secret); err != nil {
		return err
	}

	return nil
}

func (vc *VaultSystemClient) unseal(vaultKeys *VaultKeys) error {
	_, err := vc.vaultClient.Sys().UnsealWithOptions(
		&vaultapi.UnsealOpts{
			Key: vaultKeys.UnsealKeys[0],
		})
	if err != nil {
		return err
	}

	return nil
}

func (vc *VaultSystemClient) registerPlugins(plugins []*vaultapi.RegisterPluginInput) error {
	for _, plugin := range plugins {
		err := vc.vaultClient.Sys().RegisterPlugin(plugin)
		if err != nil {
			return err
		}
	}

	return nil
}

func (vc *VaultSystemClient) enableKubernetesAuth(ctx context.Context, key types.NamespacedName) error {
	authEnabled, err := vc.isKubernetesAuthEnabled()
	if err != nil {
		return err
	}

	if !authEnabled {
		err = vc.vaultClient.Sys().EnableAuthWithOptions(VaultAuthKubernetesPath,
			&vaultapi.EnableAuthOptions{
				Type: VaultAuthKubernetesType,
			})
		if err != nil {
			return err
		}
	}

	return nil
}

func (vc *VaultSystemClient) configureKubernetesAuth(ctx context.Context, key types.NamespacedName) error {
	caData, err := vc.getKubernetesAuthConfig(ctx, key)
	if err != nil {
		return err
	}

	payload := map[string]interface{}{
		VaultKubernetesHost:   DefaultVaultKubernetesHost,
		VaultKubernetesCACert: caData.CACertificate,
		VaultTokenReviewerJWT: caData.Token,
	}

	if _, err := vc.vaultClient.Logical().Write(VaultAuthKubernetesConfig, payload); err != nil {
		return err
	}

	return nil
}

func (vc *VaultSystemClient) isKubernetesAuthEnabled() (bool, error) {
	auths, err := vc.vaultClient.Logical().Read(VaultSysAuth)
	if err != nil {
		return false, err
	}

	if k8sAuth, ok := auths.Data[VaultAuthKubernetesPath]; k8sAuth != nil && ok {
		return true, nil
	}

	return false, nil
}

func (vc *VaultSystemClient) getKubernetesAuthConfig(ctx context.Context, key types.NamespacedName) (*kubernetes.KubernetesCertificateData, error) {
	vsa := &corev1.ServiceAccount{}
	err := vc.kubeClient.Get(ctx, key, vsa)
	if err != nil {
		return nil, err
	}

	vaultSecret := &corev1.Secret{}
	vaultSecretKey := types.NamespacedName{
		Name:      vsa.Secrets[0].Name,
		Namespace: key.Namespace,
	}

	err = vc.kubeClient.Get(ctx, vaultSecretKey, vaultSecret)
	if err != nil {
		return nil, err
	}

	ca := string(vaultSecret.Data[kubernetes.KubernetesSecretDataCACertificate])
	token := string(vaultSecret.Data[kubernetes.KubernetesSecretDataToken])

	return &kubernetes.KubernetesCertificateData{
		CACertificate: ca,
		Token:         token,
	}, nil
}

func NewVaultClient(vaultClient *vaultapi.Client, kubeClient client.Client) *VaultSystemClient {
	return &VaultSystemClient{
		kubeClient:  kubeClient,
		vaultClient: vaultClient,
	}
}
