package dev

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	installerv1alpha1 "github.com/puppetlabs/relay-core/pkg/apis/install.relay.sh/v1alpha1"
	"github.com/puppetlabs/relay/pkg/kubernetes"
	"github.com/puppetlabs/relay/pkg/vault"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	vaultIdentifier = "vault"
	vaultImage      = "vault:1.8.2"
	vaultAddr       = "http://vault:8200"

	vaultInitImage = "relaysh/relay-vault-init"
)

var vaultWriteValuesJobTTL = int32(120)

type vaultManagerObjects struct {
	credentialsSecret  corev1.Secret
	serviceAccount     corev1.ServiceAccount
	clusterRoleBinding rbacv1.ClusterRoleBinding
}

func newVaultManagerObjects() *vaultManagerObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      vaultIdentifier,
		Namespace: systemNamespace,
	}

	return &vaultManagerObjects{
		credentialsSecret:  corev1.Secret{ObjectMeta: objectMeta},
		serviceAccount:     corev1.ServiceAccount{ObjectMeta: objectMeta},
		clusterRoleBinding: rbacv1.ClusterRoleBinding{ObjectMeta: objectMeta},
	}
}

type vaultManager struct {
	cl      *Client
	objects *vaultManagerObjects

	cfg Config
}

func (m *vaultManager) reconcile(ctx context.Context) error {
	if err := m.reconcileInit(ctx); err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) reconcileInit(ctx context.Context) error {
	cl := m.cl.APIClient

	saKey := client.ObjectKeyFromObject(&m.objects.serviceAccount)

	if err := cl.Get(ctx, saKey, &m.objects.serviceAccount); err != nil {
		return err
	}

	m.rbacDefinition(&m.objects.clusterRoleBinding)
	if err := cl.Create(ctx, &m.objects.clusterRoleBinding); err != nil && !k8serrors.IsAlreadyExists(err) {
		return err
	}

	initJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-init", vaultIdentifier),
			Namespace: systemNamespace,
		},
	}
	err := m.getJob(ctx, initJob)
	if err != nil && k8serrors.IsNotFound(err) {
		initJob := m.initJob()

		if err := cl.Create(ctx, initJob); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err := m.waitForJobCompletion(ctx, initJob); err != nil {
		return err
	}

	if err := m.cleanupJobs(ctx, []*batchv1.Job{initJob}); err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) rbacDefinition(crb *rbacv1.ClusterRoleBinding) {
	crb.RoleRef =
		rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "cluster-admin",
		}
	crb.Subjects = []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      "default",
			Namespace: systemNamespace,
		},
	}
}

func (m *vaultManager) initJob() *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-init", vaultIdentifier),
			Namespace: systemNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "vault-action",
							Image: vaultInitImage,
							Env: []corev1.EnvVar{
								{Name: "VAULT_ADDR", Value: vaultAddr},
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
		},
	}
}

func (m *vaultManager) baseJob(job *batchv1.Job) {
	job.Spec = batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyNever,
			},
		},
	}
}

func (m *vaultManager) baseJobContainer(container *corev1.Container) {
	container.Name = "vault-action"
	container.Image = vaultImage
	container.Env = []corev1.EnvVar{
		{Name: "VAULT_ADDR", Value: vaultAddr},
	}
}

func (m *vaultManager) credentialsEnvs(container *corev1.Container) {
	envs := []corev1.EnvVar{
		{
			Name: "VAULT_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "root-token",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.objects.credentialsSecret.Name,
					},
				},
			},
		},
		{
			Name: "VAULT_UNSEAL_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "unseal-key",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: m.objects.credentialsSecret.Name,
					},
				},
			},
		},
	}

	container.Env = append(container.Env, envs...)
}

func (m *vaultManager) relayCoreAccessJob(rc *installerv1alpha1.RelayCore, job *batchv1.Job, configMap *corev1.ConfigMap) error {
	policyGen := newVaultPolicyGenerator(rc)

	tenantConfig, err := policyGen.metadataAPITenantConfigFile()
	if err != nil {
		return err
	}

	configMap.Data = map[string]string{
		"log-service.hcl":                 string(policyGen.logServiceFile()),
		"operator.hcl":                    string(policyGen.operatorFile()),
		"metadata-api.hcl":                string(policyGen.metadataAPIFile()),
		"metadata-api-tenant.hcl":         string(policyGen.metadataAPITenantFile()),
		"metadata-api-tenant-config.json": string(tenantConfig),
		"config.sh":                       vaultAccessScript,
	}

	m.baseJob(job)

	container := corev1.Container{}

	m.baseJobContainer(&container)
	m.credentialsEnvs(&container)
	m.relayCoreAccessJobEnvs(rc, &container)

	container.VolumeMounts = []corev1.VolumeMount{
		{
			Name:      "policy-config",
			MountPath: "/vault-policy-config",
		},
	}

	container.Command = []string{"/bin/sh", "/vault-policy-config/config.sh"}

	job.Spec.Template.Spec.Containers = []corev1.Container{container}
	job.Spec.TTLSecondsAfterFinished = &vaultWriteValuesJobTTL

	job.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "policy-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: configMap.Name,
					},
				},
			},
		},
	}

	return nil
}

func (m *vaultManager) relayCoreAccessJobEnvs(rc *installerv1alpha1.RelayCore, container *corev1.Container) {
	container.Env = append(container.Env,
		corev1.EnvVar{
			Name: "JWT_SIGNING_PUBLIC_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "public-key.pem",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: rc.Status.Vault.JWTSigningKeySecret,
					},
				},
			},
		},
		corev1.EnvVar{
			Name:  "LOG_SERVICE_POLICY",
			Value: rc.Status.Vault.LogServiceRole,
		},
		corev1.EnvVar{
			Name:  "LOG_SERVICE_SERVICE_ACCOUNT_NAME",
			Value: rc.Status.Vault.LogServiceServiceAccount,
		},
		corev1.EnvVar{
			Name:  "OPERATOR_POLICY",
			Value: rc.Status.Vault.OperatorRole,
		},
		corev1.EnvVar{
			Name:  "OPERATOR_SERVICE_ACCOUNT_NAME",
			Value: rc.Status.Vault.OperatorServiceAccount,
		},
		corev1.EnvVar{
			Name:  "METADATA_API_POLICY",
			Value: rc.Status.Vault.MetadataAPIRole,
		},
		corev1.EnvVar{
			Name:  "METADATA_API_TENANT_POLICY",
			Value: fmt.Sprintf("%s-tenant", rc.Status.Vault.MetadataAPIRole),
		},
		corev1.EnvVar{
			Name:  "METADATA_API_SERVICE_ACCOUNT_NAME",
			Value: rc.Status.Vault.MetadataAPIServiceAccount,
		},
		corev1.EnvVar{
			Name:  "SERVICE_ACCOUNT_NAMESPACE",
			Value: rc.Namespace,
		},
		corev1.EnvVar{
			Name:  "JWT_AUTH_ROLE",
			Value: rc.Spec.MetadataAPI.VaultAuthRole,
		},
		corev1.EnvVar{
			Name:  "JWT_AUTH_PATH",
			Value: rc.Spec.MetadataAPI.VaultAuthPath,
		},
		corev1.EnvVar{
			Name:  "LOG_SERVICE_PATH",
			Value: rc.Spec.Vault.LogServicePath,
		},
		corev1.EnvVar{
			Name:  "TENANT_PATH",
			Value: rc.Spec.Vault.TenantPath,
		},
		corev1.EnvVar{
			Name:  "TRANSIT_PATH",
			Value: rc.Spec.Vault.TransitPath,
		},
		corev1.EnvVar{
			Name:  "TRANSIT_KEY",
			Value: rc.Spec.Vault.TransitKey,
		},
	)
}

func (m *vaultManager) writeValuesJob(vals map[string]string, job *batchv1.Job) {
	m.baseJob(job)

	container := corev1.Container{}

	m.baseJobContainer(&container)
	m.credentialsEnvs(&container)

	cmds := []string{}

	for k, v := range vals {
		cmds = append(cmds, fmt.Sprintf("vault kv put %s value='%s'", k, v))
	}

	script := strings.Join(cmds, "; ")

	container.Command = []string{"/bin/sh", "-c", script}

	job.Spec.Template.Spec.Containers = []corev1.Container{container}
	job.Spec.TTLSecondsAfterFinished = &vaultWriteValuesJobTTL
}

func (m *vaultManager) addRelayCoreAccess(ctx context.Context, rc *installerv1alpha1.RelayCore) error {
	objectMeta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-vault-access", rc.Name),
		Namespace: systemNamespace,
	}

	job := batchv1.Job{ObjectMeta: objectMeta}
	configMap := corev1.ConfigMap{ObjectMeta: objectMeta}

	if err := m.relayCoreAccessJob(rc, &job, &configMap); err != nil {
		return err
	}

	if err := m.cl.APIClient.Create(ctx, &job); err != nil {
		return err
	}

	if err := m.cl.APIClient.Create(ctx, &configMap); err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) writeSecrets(ctx context.Context, vals map[string]string) error {
	job := batchv1.Job{ObjectMeta: metav1.ObjectMeta{
		GenerateName: "vault-write-values-",
		Namespace:    systemNamespace,
	}}

	m.writeValuesJob(vals, &job)

	if err := m.cl.APIClient.Create(ctx, &job); err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) getJob(ctx context.Context, job *batchv1.Job) error {
	cl := m.cl.APIClient

	key := client.ObjectKeyFromObject(job)

	return cl.Get(ctx, key, job)
}

func (m *vaultManager) waitForJobCompletion(ctx context.Context, job *batchv1.Job) error {
	err := retry.Wait(ctx, func(ctx context.Context) (bool, error) {
		if err := m.getJob(ctx, job); err != nil {
			return retry.Repeat(err)
		}

		if len(job.Status.Conditions) == 0 {
			return retry.Repeat(errors.New("waiting for vault operator job to finish"))
		}

		for _, cond := range job.Status.Conditions {
			switch cond.Type {
			case batchv1.JobComplete:
				if cond.Status == corev1.ConditionTrue {
					return retry.Done(nil)
				}
			case batchv1.JobFailed:
				if cond.Status == corev1.ConditionTrue {
					return retry.Done(errors.New(cond.Message))
				}
			}
		}

		return retry.Repeat(nil)
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) cleanupJobs(ctx context.Context, jobs []*batchv1.Job) error {
	for _, job := range jobs {
		policy := client.PropagationPolicy(metav1.DeletePropagationBackground)

		if err := m.cl.APIClient.Delete(ctx, job, policy); err != nil {
			return err
		}
	}

	return nil
}

func newVaultManager(cl *Client, cfg Config) *vaultManager {
	return &vaultManager{
		cl:      cl,
		objects: newVaultManagerObjects(),
		cfg:     cfg,
	}
}

// vaultPolicy is a subset of the vault.Policy model. We created this here for
// 2 reasons: vault is a big module that we want to prevent pulling in as much
// as possible (the sdk and api are usually okay) and the model isn't made for
// encoding, so it doesn't include the proper hcl tags.
type vaultPolicy struct {
	Paths []vaultPolicyPath `hcl:"path,block"`
}

type vaultPolicyPath struct {
	Name         string   `hcl:"name,label"`
	Capabilities []string `hcl:"capabilities"`
}

type vaultPolicyGenerator struct {
	rc *installerv1alpha1.RelayCore
}

func (g *vaultPolicyGenerator) operatorFile() []byte {
	policy := vaultPolicy{
		Paths: []vaultPolicyPath{
			{
				Name:         path.Join(g.rc.Spec.Vault.TransitPath, "encrypt", g.rc.Spec.Vault.TransitKey),
				Capabilities: []string{"update"},
			},
			{
				Name:         path.Join("sys", "mounts"),
				Capabilities: []string{"read"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.TenantPath, "data", "workflows", "*"),
				Capabilities: []string{"create", "update"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.TenantPath, "data", "connections", "*"),
				Capabilities: []string{"create", "update"},
			},
			{
				Name:         path.Join("oauth", "+", "creds", "*"),
				Capabilities: []string{"read", "create", "update", "delete"},
			},
			{
				Name:         path.Join("oauth", "auth0-management", "self", "management-api"),
				Capabilities: []string{"read"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.TenantPath, "metadata", "workflows", "*"),
				Capabilities: []string{"list", "delete"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.TenantPath, "metadata", "connections", "*"),
				Capabilities: []string{"list", "delete"},
			},
			{
				Name:         path.Join("customers-oauth", "auth-code-url"),
				Capabilities: []string{"update"},
			},
			{
				Name:         path.Join("customers-oauth/creds/*"),
				Capabilities: []string{"read", "create", "update", "delete"},
			},
		},
	}

	return g.generate(&policy)
}

func (g *vaultPolicyGenerator) logServiceFile() []byte {
	policy := vaultPolicy{
		Paths: []vaultPolicyPath{
			{
				Name:         path.Join("sys", "mounts"),
				Capabilities: []string{"read"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.LogServicePath, "data", "logs", "*"),
				Capabilities: []string{"create", "read", "update"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.LogServicePath, "data", "contexts", "*"),
				Capabilities: []string{"create", "read", "update"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.LogServicePath, "metadata", "logs", "*"),
				Capabilities: []string{"list", "delete"},
			},
			{
				Name:         path.Join(g.rc.Spec.Vault.LogServicePath, "metadata", "contexts", "*"),
				Capabilities: []string{"list", "delete"},
			},
		},
	}

	return g.generate(&policy)
}

func (g *vaultPolicyGenerator) metadataAPIFile() []byte {
	policy := vaultPolicy{
		Paths: []vaultPolicyPath{
			{
				Name:         path.Join(g.rc.Spec.Vault.TransitPath, "decrypt", g.rc.Spec.Vault.TransitKey),
				Capabilities: []string{"update"},
			},
		},
	}

	return g.generate(&policy)
}

func (g *vaultPolicyGenerator) metadataAPITenantFile() []byte {
	spec := g.rc.Spec

	tenantEntity := "{{identity.entity.aliases.$AUTH_JWT_ACCESSOR.metadata.tenant_id}}"
	domainEntity := "{{identity.entity.aliases.$AUTH_JWT_ACCESSOR.metadata.domain_id}}"

	policy := vaultPolicy{
		Paths: []vaultPolicyPath{
			{
				Name:         path.Join(spec.Vault.TenantPath, "data", "workflows", tenantEntity, "*"),
				Capabilities: []string{"read"},
			},
			{
				Name:         path.Join(spec.Vault.TenantPath, "metadata", "connections", domainEntity, "*"),
				Capabilities: []string{"list"},
			},
			{
				Name:         path.Join(spec.Vault.TenantPath, "data", "connections", domainEntity, "*"),
				Capabilities: []string{"read"},
			},
		},
	}

	return g.generate(&policy)
}

func (g *vaultPolicyGenerator) metadataAPITenantConfigFile() ([]byte, error) {
	file := struct {
		Name           string            `json:"name"`
		RoleType       string            `json:"role_type"`
		BoundAudiences []string          `json:"bound_audiences"`
		UserClaim      string            `json:"user_claim"`
		TokenType      string            `json:"token_type"`
		TokenPolicies  []string          `json:"token_policies"`
		ClaimMappings  map[string]string `json:"claim_mappings"`
	}{
		Name:           g.rc.Spec.MetadataAPI.VaultAuthRole,
		RoleType:       "jwt",
		BoundAudiences: []string{"k8s.relay.sh/metadata-api/v1"},
		UserClaim:      "sub",
		TokenType:      "batch",
		TokenPolicies:  []string{fmt.Sprintf("%s-tenant", g.rc.Status.Vault.MetadataAPIRole)},
		ClaimMappings: map[string]string{
			"relay.sh/domain-id": "domain_id",
			"relay.sh/tenant-id": "tenant_id",
		},
	}

	return json.Marshal(file)
}

func (g *vaultPolicyGenerator) generate(policy *vaultPolicy) []byte {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(policy, f.Body())

	return f.Bytes()
}

func newVaultPolicyGenerator(rc *installerv1alpha1.RelayCore) *vaultPolicyGenerator {
	return &vaultPolicyGenerator{rc: rc}
}

var (
	vaultAccessScript = `
apk add --no-cache jq gettext
vault auth enable -path=jwt-tenants jwt
vault write ${JWT_AUTH_PATH}/config \
   jwt_supported_algs="RS256,RS512" \
   jwt_validation_pubkeys="${JWT_SIGNING_PUBLIC_KEY}"
export AUTH_JWT_ACCESSOR="$( vault auth list -format=json | jq -r '."jwt-tenants/".accessor' )"
vault secrets enable -path=${LOG_SERVICE_PATH} kv-v2
vault secrets enable -path=${TENANT_PATH} kv-v2
vault secrets enable -path=${TRANSIT_PATH} transit
vault write ${TRANSIT_PATH}/keys/${TRANSIT_KEY} derived=true
vault policy write ${OPERATOR_POLICY} /vault-policy-config/operator.hcl
vault policy write ${METADATA_API_POLICY} /vault-policy-config/metadata-api.hcl
vault policy write ${LOG_SERVICE_POLICY} /vault-policy-config/log-service.hcl
envsubst '$AUTH_JWT_ACCESSOR' < /vault-policy-config/metadata-api-tenant.hcl > /tmp/metadata-api-tenant.hcl
vault policy write ${METADATA_API_TENANT_POLICY} /tmp/metadata-api-tenant.hcl
vault write auth/kubernetes/role/${OPERATOR_POLICY} \
    bound_service_account_names=${OPERATOR_SERVICE_ACCOUNT_NAME} \
    bound_service_account_namespaces=${SERVICE_ACCOUNT_NAMESPACE} \
    ttl=24h policies=${OPERATOR_POLICY}
vault write auth/kubernetes/role/${METADATA_API_POLICY} \
    bound_service_account_names=${METADATA_API_SERVICE_ACCOUNT_NAME} \
    bound_service_account_namespaces=${SERVICE_ACCOUNT_NAMESPACE} \
	ttl=24h policies=${METADATA_API_POLICY}
vault write auth/kubernetes/role/${LOG_SERVICE_POLICY} \
    bound_service_account_names=${LOG_SERVICE_SERVICE_ACCOUNT_NAME} \
    bound_service_account_namespaces=${SERVICE_ACCOUNT_NAMESPACE} \
    ttl=24h \
    policies=${LOG_SERVICE_POLICY}
vault write ${JWT_AUTH_PATH}/role/${JWT_AUTH_ROLE} - < /vault-policy-config/metadata-api-tenant-config.json
`
)

type vaultInitializer struct {
	client *vault.VaultSystemClient
}

func (vi *vaultInitializer) Initialize(ctx context.Context) error {
	plugins := []*vaultapi.RegisterPluginInput{
		{
			Name:    "oauthapp",
			Type:    consts.PluginTypeSecrets,
			SHA256:  "2455ad9450e415efb42fd28c1436f4bf7e377524be5e534e55e6658b8ef56bd2",
			Command: "vault-plugin-secrets-oauthapp-v3.0.0-beta.3-linux-amd64",
		},
	}

	key := types.NamespacedName{
		Name:      vaultIdentifier,
		Namespace: systemNamespace,
	}

	if err := vi.client.InitializeVault(ctx, key, plugins); err != nil {
		return err
	}

	return nil
}

func NewVaultInitializer(vaultClient *vaultapi.Client) (*vaultInitializer, error) {
	kubeClient, err := kubernetes.NewKubeClient(DefaultScheme)
	if err != nil {
		return nil, err
	}

	return &vaultInitializer{
		client: vault.NewVaultClient(vaultClient, kubeClient),
	}, nil
}
