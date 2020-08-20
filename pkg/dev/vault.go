package dev

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/puppetlabs/relay-core/pkg/util/retry"
	"github.com/puppetlabs/relay/pkg/cluster"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	vaultImage                  = "vault:1.5.0"
	vaultCredentialsSecretName  = "vault-credentials"
	vaultInitJobName            = "vault-init"
	vaultUnsealJobName          = "vault-unseal"
	vaultConfigureJobName       = "vault-configure"
	vaultCredentialsStorageName = "vault-credential-storage"
	vaultInitMountPath          = "/vault-init"
	vaultInitDataFile           = "init-data.json"
	vaultAddr                   = "http://vault:8200"
)

type vaultKeys struct {
	UnsealKeys []string `json:"unseal_keys_b64"`
	RootToken  string   `json:"root_token"`
}

type vaultManager struct {
	cl        *cluster.Client
	namespace string

	cfg Config
}

func (m *vaultManager) init(ctx context.Context) error {
	key := client.ObjectKey{
		Name:      vaultInitJobName,
		Namespace: "relay-system",
	}

	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Mi"),
				},
			},
		},
	}

	if err := m.cl.APIClient.Create(ctx, pvc); err != nil {
		return err
	}

	initDataPath := filepath.Join(vaultInitMountPath, vaultInitDataFile)
	cmds := []string{
		"/bin/sh",
		"-c",
		fmt.Sprintf("vault operator init -format=json -key-shares=1 -key-threshold=1 > %s", initDataPath),
	}

	job := vaultOperatorJobWithInitVolume(key, cmds)

	if err := runJobAndWait(ctx, m.cl, job); err != nil {
		return err
	}

	pvcKey, err := client.ObjectKeyFromObject(pvc)
	if err != nil {
		return err
	}

	if err := m.cl.APIClient.Get(ctx, pvcKey, pvc); err != nil {
		return err
	}

	localInitDataPath := filepath.Join(
		m.cfg.DataDir,
		cluster.HostStorageName,
		pvc.Spec.VolumeName,
		vaultInitDataFile)
	bytes, err := ioutil.ReadFile(localInitDataPath)
	if err != nil {
		return err
	}

	var keys vaultKeys

	if err := json.Unmarshal(bytes, &keys); err != nil {
		return err
	}

	sec := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      vaultCredentialsSecretName,
			Namespace: "relay-system",
		},
		StringData: map[string]string{
			"root-token": keys.RootToken,
			"unseal-key": keys.UnsealKeys[0],
		},
	}

	if err := m.cl.APIClient.Create(ctx, sec); err != nil {
		return err
	}

	if err := m.cl.APIClient.Delete(ctx, pvc); err != nil {
		return err
	}

	if err := m.unseal(ctx); err != nil {
		return err
	}

	return m.configure(ctx)
}

func (m *vaultManager) configure(ctx context.Context) error {
	key := client.ObjectKey{
		Name:      vaultConfigureJobName,
		Namespace: "relay-system",
	}

	script := `
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
EOT
vault write auth/kubernetes/role/relay-metadata-api \
    bound_service_account_names=relay-metadata-api-vault \
    bound_service_account_namespaces=relay-system \
    ttl=24h \
    policies=relay/metadata-api
vault write auth/kubernetes/role/relay-tasks \
    bound_service_account_names=relay-tasks-vault \
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

	sa := &corev1.ServiceAccount{}
	saKey := client.ObjectKey{
		Name:      "vault",
		Namespace: "relay-system",
	}

	if err := m.cl.APIClient.Get(ctx, saKey, sa); err != nil {
		return err
	}

	// I think we can assume there's only 1???
	saSecret := sa.Secrets[0]

	envs := []corev1.EnvVar{
		{
			Name: "VAULT_CA_CERT",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "ca.crt",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: saSecret.Name,
					},
				},
			},
		},
		{
			Name: "VAULT_JWT_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "token",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: saSecret.Name,
					},
				},
			},
		},
		{
			Name: "VAULT_JWT_SINGING_PUBLIC_KEY",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "public-key.pem",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: jwtSigningKeysSecretName,
					},
				},
			},
		},
	}
	cmds := []string{"/bin/sh", "-c", script}
	job := vaultOperatorJobWithAuth(key, cmds, envs...)

	return runJobAndWait(ctx, m.cl, job)
}

func (m *vaultManager) unseal(ctx context.Context) error {
	key := client.ObjectKey{
		Name:      vaultUnsealJobName,
		Namespace: "relay-system",
	}

	cmds := []string{
		"/bin/sh",
		"-c",
		"vault operator unseal ${VAULT_UNSEAL_KEY}",
	}
	job := vaultOperatorJobWithAuth(key, cmds)

	return runJobAndWait(ctx, m.cl, job)
}

func newVaultManager(cl *cluster.Client, cfg Config) *vaultManager {
	return &vaultManager{
		cl:  cl,
		cfg: cfg,
	}
}

func vaultOperatorJob(key client.ObjectKey, cmds []string) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      key.Name,
			Namespace: key.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    key.Name,
							Image:   vaultImage,
							Command: cmds,
							Env: []corev1.EnvVar{
								{Name: "VAULT_ADDR", Value: vaultAddr},
							},
						},
					},
				},
			},
		},
	}
}

func vaultOperatorJobWithInitVolume(key client.ObjectKey, cmds []string) *batchv1.Job {
	job := vaultOperatorJob(key, cmds)

	initVolume := corev1.Volume{
		Name: vaultCredentialsStorageName,
		VolumeSource: corev1.VolumeSource{
			PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
				ClaimName: key.Name,
			},
		},
	}

	volumes := []corev1.Volume{initVolume}
	volumeMounts := []corev1.VolumeMount{
		{Name: initVolume.Name, MountPath: vaultInitMountPath},
	}

	job.Spec.Template.Spec.Volumes = volumes

	for i := range job.Spec.Template.Spec.Containers {
		job.Spec.Template.Spec.Containers[i].VolumeMounts = volumeMounts
	}

	return job
}

func vaultOperatorJobWithAuth(key client.ObjectKey, cmds []string, envs ...corev1.EnvVar) *batchv1.Job {
	job := vaultOperatorJob(key, cmds)

	creds := []corev1.EnvVar{
		{
			Name: "VAULT_TOKEN",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "root-token",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: vaultCredentialsSecretName,
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
						Name: vaultCredentialsSecretName,
					},
				},
			},
		},
	}

	envs = append(creds, envs...)

	for i, container := range job.Spec.Template.Spec.Containers {
		oldEnv := container.Env

		newEnv := append(oldEnv, envs...)

		job.Spec.Template.Spec.Containers[i].Env = newEnv
	}

	return job
}

func waitForJobToComplete(ctx context.Context, cl *cluster.Client, job *batchv1.Job) error {
	err := retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
		key, err := client.ObjectKeyFromObject(job)
		if err != nil {
			return retry.RetryPermanent(err)
		}

		if err := cl.APIClient.Get(ctx, key, job); err != nil {
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

// runJobAndWait runs a given job and then watches its conditions for
// completion and returns
//
// TODO move this to a common package or file if we end up using more batch
// jobs.
func runJobAndWait(ctx context.Context, cl *cluster.Client, job *batchv1.Job) error {
	if err := cl.APIClient.Create(ctx, job); err != nil {
		return err
	}

	return waitForJobToComplete(ctx, cl, job)
}
