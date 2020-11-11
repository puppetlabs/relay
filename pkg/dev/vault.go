package dev

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

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

	jwtSigningKeyEnv := secretEnv.DeepCopy()
	jwtSigningKeyEnv.Name = "VAULT_JWT_SINGING_PUBLIC_KEY"
	jwtSigningKeyEnv.ValueFrom.SecretKeyRef.Key = "public-key.pem"
	// TODO I think we can just get this from the RelayCore object
	jwtSigningKeyEnv.ValueFrom.SecretKeyRef.LocalObjectReference.Name = jwtSigningKeysSecretName

	container.Env = append(container.Env, *certEnv, *jwtTokenEnv, *jwtSigningKeyEnv)

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

	container.Command = []string{script}

	job.Spec.Template.Spec.Containers = []corev1.Container{container}
	job.Spec.TTLSecondsAfterFinished = &vaultWriteValuesJobTTL
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

	err = retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
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
					return retry.RetryPermanent(errors.New(cond.Message))
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

var vaultConfigureScript = `
apk add --no-cache jq
vault auth enable -path=kubernetes kubernetes
vault write auth/kubernetes/config \
    kubernetes_host="https://kubernetes.default.svc" \
    kubernetes_ca_cert="${VAULT_CA_CERT}" \
    token_reviewer_jwt="${VAULT_JWT_TOKEN}"
vault auth enable -path=jwt-tenants jwt
vault write auth/jwt-tenants/config \
   jwt_supported_algs="RS256,RS512" \
   jwt_validation_pubkeys="${VAULT_JWT_SINGING_PUBLIC_KEY}"
AUTH_JWT_ACCESSOR="$( vault auth list -format=json | jq -r '."jwt-tenants/".accessor' )"
vault secrets enable -path=customers kv-v2
vault secrets enable -path=transit-tenants transit
vault write transit-tenants/keys/metadata-api derived=true
vault policy write relay/metadata-api - <<EOT
path "transit-tenants/decrypt/metadata-api" {
  capabilities = ["update"]
}
EOT
vault policy write relay/metadata-api-tenant - <<EOT
path "customers/data/workflows/{{identity.entity.aliases.${AUTH_JWT_ACCESSOR}.metadata.tenant_id}}/*" {
    capabilities = ["read"]
}
path "customers/metadata/connections/{{identity.entity.aliases.${AUTH_JWT_ACCESSOR}.metadata.domain_id}}/*" {
    capabilities = ["list"]
}
path "customers/data/connections/{{identity.entity.aliases.${AUTH_JWT_ACCESSOR}.metadata.domain_id}}/*" {
    capabilities = ["read"]
}
EOT
vault policy write relay/tasks - <<EOT
path "transit-tenants/encrypt/metadata-api" {
  capabilities = ["update"]
}
path "sys/mounts" {
  capabilities = ["read"]
}
path "customers/data/workflows/*" {
  capabilities = ["create", "update"]
}
path "customers/data/connections/*" {
  capabilities = ["create", "update"]
}
path "customers/metadata/workflows/*" {
  capabilities = ["list", "delete"]
}
path "customers/metadata/connections/*" {
  capabilities = ["list", "delete"]
}
path "oauth/+/config/auth_code_url" {
  capabilities = ["update"]
}
path "oauth/+/creds/*" {
  capabilities = ["read", "create", "update", "delete"]
}
EOT
vault write auth/kubernetes/role/relay-core-v1-metadata-api \
    bound_service_account_names=relay-core-v1-metadata-api-vault \
    bound_service_account_namespaces=relay-system \
    ttl=24h \
    policies=relay/metadata-api
vault write auth/kubernetes/role/relay-core-v1-operator \
    bound_service_account_names=relay-core-v1-operator-vault \
    bound_service_account_namespaces=relay-system \
    ttl=24h \
    policies=relay/tasks
vault write auth/jwt-tenants/role/tenant - <<EOT
{
    "name": "tenant",
    "role_type": "jwt",
    "bound_audiences": ["k8s.relay.sh/metadata-api/v1"],
    "user_claim": "sub",
    "token_type": "batch",
    "token_policies": ["relay/metadata-api-tenant"],
    "claim_mappings": {
        "relay.sh/domain-id": "domain_id",
        "relay.sh/tenant-id": "tenant_id"
    }
}
EOT
`
