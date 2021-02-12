package dev

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hclwrite"
	installerv1alpha1 "github.com/puppetlabs/relay-core/pkg/apis/install.relay.sh/v1alpha1"
	"github.com/puppetlabs/relay-core/pkg/util/retry"
	"github.com/puppetlabs/relay/pkg/cluster"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	vaultImage                 = "vault:1.5.0"
	vaultAddr                  = "http://vault:8200"
	vaultInitResultStorageName = "vault-init-result-storage"
	vaultInitResultDataFile    = "init-data.json"
	vaultInitResultMountPath   = "/vault-init"
	vaultInitVolumeSize        = "1Mi"
)

var vaultWriteValuesJobTTL = int32(120)

type vaultManagerObjects struct {
	initPVC           corev1.PersistentVolumeClaim
	initJob           batchv1.Job
	configureJob      batchv1.Job
	unsealJob         batchv1.Job
	credentialsSecret corev1.Secret
	serviceAccount    corev1.ServiceAccount
}

func newVaultManagerObjects() *vaultManagerObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      "vault",
		Namespace: systemNamespace,
	}

	return &vaultManagerObjects{
		initPVC: corev1.PersistentVolumeClaim{ObjectMeta: objectMeta},
		initJob: batchv1.Job{ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-init", objectMeta.Name),
			Namespace: objectMeta.Namespace,
		}},
		configureJob: batchv1.Job{ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-configure-init", objectMeta.Name),
			Namespace: objectMeta.Namespace,
		}},
		unsealJob: batchv1.Job{ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-unseal-init", objectMeta.Name),
			Namespace: objectMeta.Namespace,
		}},
		credentialsSecret: corev1.Secret{ObjectMeta: objectMeta},
		serviceAccount:    corev1.ServiceAccount{ObjectMeta: objectMeta},
	}
}

type vaultKeys struct {
	UnsealKeys []string `json:"unseal_keys_b64"`
	RootToken  string   `json:"root_token"`
}

type vaultManager struct {
	cl      *cluster.Client
	objects *vaultManagerObjects

	cfg Config
}

func (m *vaultManager) reconcile(ctx context.Context) error {
	if err := m.reconcileInit(ctx); err != nil {
		return err
	}

	if err := m.reconcileUnseal(ctx); err != nil {
		return err
	}

	if err := m.reconcileConfiguration(ctx); err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) reconcileInit(ctx context.Context) error {
	cl := m.cl.APIClient

	saKey, err := client.ObjectKeyFromObject(&m.objects.serviceAccount)
	if err != nil {
		return err
	}

	if err := cl.Get(ctx, saKey, &m.objects.serviceAccount); err != nil {
		return err
	}

	err = m.waitForJobCompletion(ctx, &m.objects.initJob)

	if err != nil && k8serrors.IsNotFound(err) {
		m.initPVC(&m.objects.initPVC)
		if err := cl.Create(ctx, &m.objects.initPVC); err != nil {
			return err
		}

		m.initJob(&m.objects.initJob)
		if err := cl.Create(ctx, &m.objects.initJob); err != nil {
			return err
		}

		if err := m.waitForJobCompletion(ctx, &m.objects.initJob); err != nil {
			return err
		}

		credentials, err := m.decodeCredentialsFromInitResult(ctx)
		if err != nil {
			return err
		}

		m.credentialsSecret(*credentials, &m.objects.credentialsSecret)

		if err := cl.Create(ctx, &m.objects.credentialsSecret); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func (m *vaultManager) reconcileConfiguration(ctx context.Context) error {
	err := m.waitForJobCompletion(ctx, &m.objects.configureJob)

	if err != nil && k8serrors.IsNotFound(err) {
		m.configureJob(&m.objects.configureJob)
		if err := m.cl.APIClient.Create(ctx, &m.objects.configureJob); err != nil {
			return err
		}

		if err := m.waitForJobCompletion(ctx, &m.objects.configureJob); err != nil {
			return err
		}
	}

	return nil
}

func (m *vaultManager) reconcileUnseal(ctx context.Context) error {
	err := m.waitForJobCompletion(ctx, &m.objects.unsealJob)

	if err != nil && k8serrors.IsNotFound(err) {
		ss := appsv1.StatefulSet{}
		ssKey := client.ObjectKey{Name: "vault", Namespace: systemNamespace}

		if err := m.cl.APIClient.Get(ctx, ssKey, &ss); err != nil {
			return err
		}

		if ss.Status.ReadyReplicas != *ss.Spec.Replicas {
			m.unsealJob(&m.objects.unsealJob)
			if err := m.cl.APIClient.Create(ctx, &m.objects.unsealJob); err != nil {
				return err
			}

			if err := m.waitForJobCompletion(ctx, &m.objects.unsealJob); err != nil {
				return err
			}

			if err := m.cleanupJobs(ctx, []*batchv1.Job{&m.objects.unsealJob}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *vaultManager) initPVC(pvc *corev1.PersistentVolumeClaim) {
	pvc.Spec = corev1.PersistentVolumeClaimSpec{
		AccessModes: []corev1.PersistentVolumeAccessMode{
			corev1.ReadWriteOnce,
		},
		Resources: corev1.ResourceRequirements{
			Requests: corev1.ResourceList{
				corev1.ResourceStorage: resource.MustParse(vaultInitVolumeSize),
			},
		},
	}
}

func (m *vaultManager) initJob(job *batchv1.Job) {
	m.baseJob(job)

	container := corev1.Container{}
	m.baseJobContainer(&container)

	initResultDataPath := filepath.Join(vaultInitResultMountPath, vaultInitResultDataFile)
	cmd := fmt.Sprintf("vault operator init -format=json -key-shares=1 -key-threshold=1 > %s", initResultDataPath)
	cmds := []string{"/bin/sh", "-c", cmd}

	container.Command = cmds

	initVolume := corev1.Volume{
		Name: vaultInitResultStorageName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: m.objects.initPVC.Name,
			},
		},
	}

	job.Spec.Template.Spec.Volumes = []corev1.Volume{initVolume}

	container.VolumeMounts = []corev1.VolumeMount{
		{Name: initVolume.Name, MountPath: vaultInitResultMountPath},
	}

	job.Spec.Template.Spec.Containers = []corev1.Container{container}
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

func (m *vaultManager) credentialsSecret(keys vaultKeys, sec *corev1.Secret) {
	sec.StringData = make(map[string]string)
	sec.StringData["root-token"] = keys.RootToken
	sec.StringData["unseal-key"] = keys.UnsealKeys[0]
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

func (m *vaultManager) decodeCredentialsFromInitResult(ctx context.Context) (*vaultKeys, error) {
	var keys vaultKeys

	pvcKey, err := client.ObjectKeyFromObject(&m.objects.initPVC)
	if err != nil {
		return nil, err
	}

	pvc := corev1.PersistentVolumeClaim{}

	if err := m.cl.APIClient.Get(ctx, pvcKey, &pvc); err != nil {
		return nil, err
	}

	pvcDirName := fmt.Sprintf("%s_%s_%s", pvc.Spec.VolumeName, pvc.Namespace, pvc.Name)

	hostStorage := filepath.Join(m.cfg.WorkDir.Path, cluster.HostStorageName)
	localInitDataPath := filepath.Join(hostStorage, pvcDirName, vaultInitResultDataFile)

	bytes, err := ioutil.ReadFile(localInitDataPath)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(bytes, &keys); err != nil {
		return nil, err
	}

	return &keys, nil
}

func (m *vaultManager) configureJob(job *batchv1.Job) {
	m.baseJob(job)

	container := corev1.Container{}

	m.baseJobContainer(&container)
	m.credentialsEnvs(&container)

	container.Command = []string{"/bin/sh", "-c", vaultConfigureScript}

	secretEnv := corev1.EnvVar{
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: m.objects.serviceAccount.Secrets[0].Name,
				},
			},
		},
	}

	certEnv := secretEnv.DeepCopy()
	certEnv.Name = "VAULT_CA_CERT"
	certEnv.ValueFrom.SecretKeyRef.Key = "ca.crt"

	jwtTokenEnv := secretEnv.DeepCopy()
	jwtTokenEnv.Name = "VAULT_JWT_TOKEN"
	jwtTokenEnv.ValueFrom.SecretKeyRef.Key = "token"

	container.Env = append(container.Env, *certEnv, *jwtTokenEnv)

	job.Spec.Template.Spec.Containers = []corev1.Container{container}
}

func (m *vaultManager) unsealJob(job *batchv1.Job) {
	m.baseJob(job)

	container := corev1.Container{}

	m.baseJobContainer(&container)
	m.credentialsEnvs(&container)

	cmd := "vault operator unseal ${VAULT_UNSEAL_KEY}"
	container.Command = []string{"/bin/sh", "-c", cmd}

	job.Spec.Template.Spec.Containers = []corev1.Container{container}
	job.Spec.TTLSecondsAfterFinished = &vaultWriteValuesJobTTL
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

func (m *vaultManager) waitForJobCompletion(ctx context.Context, job *batchv1.Job) error {
	cl := m.cl.APIClient

	key, err := client.ObjectKeyFromObject(job)
	if err != nil {
		return err
	}

	var maxAttempts = 10

	err = retry.Retry(ctx, 5*time.Second, func() *retry.RetryError {
		if err := cl.Get(ctx, key, job); err != nil {
			return retry.RetryPermanent(err)
		}

		if len(job.Status.Conditions) == 0 {
			return retry.RetryTransient(errors.New("waiting for vault operator job to finish"))
		}

	conditions:
		for _, cond := range job.Status.Conditions {
			switch cond.Type {
			case batchv1.JobComplete:
				if cond.Status == corev1.ConditionTrue {
					break conditions
				}
			case batchv1.JobFailed:
				if cond.Status == corev1.ConditionTrue {
					if maxAttempts == 0 {
						return retry.RetryPermanent(errors.New(cond.Message))
					}

					maxAttempts--
				}
			}

			return retry.RetryTransient(errors.New("waiting for vault operator job to finish"))
		}

		return retry.RetryPermanent(nil)
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

func newVaultManager(cl *cluster.Client, cfg Config) *vaultManager {
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
	vaultConfigureScript = `
vault plugin register \
	-sha256=de9f96853636419150461ec34c7af9e4cf6b981f2476e1eaa00d5a58b3ddad7e \
	-command=oauthapp \
	secret oauthapp
vault auth enable -path=kubernetes kubernetes
vault write auth/kubernetes/config \
    kubernetes_host="https://kubernetes.default.svc" \
    kubernetes_ca_cert="${VAULT_CA_CERT}" \
    token_reviewer_jwt="${VAULT_JWT_TOKEN}"
`

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
