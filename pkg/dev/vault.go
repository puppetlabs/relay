package dev

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/puppetlabs/leg/timeutil/pkg/backoff"
	"github.com/puppetlabs/leg/timeutil/pkg/retry"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	vaultIdentifier = "vault"
	vaultImage      = "vault:1.8.2"
	vaultAddr       = "http://vault:8200"
)

var vaultWriteValuesJobTTL = int32(120)

type vaultManagerObjects struct {
	credentialsSecret corev1.Secret
	serviceAccount    corev1.ServiceAccount
}

func newVaultManagerObjects() *vaultManagerObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      vaultIdentifier,
		Namespace: systemNamespace,
	}

	return &vaultManagerObjects{
		credentialsSecret: corev1.Secret{ObjectMeta: objectMeta},
		serviceAccount:    corev1.ServiceAccount{ObjectMeta: objectMeta},
	}
}

type vaultManager struct {
	cl      *Client
	objects *vaultManagerObjects

	cfg Config
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

func (m *vaultManager) writeSecrets(ctx context.Context, vals map[string]string) error {
	job := batchv1.Job{ObjectMeta: metav1.ObjectMeta{
		GenerateName: "vault-write-values-",
		Namespace:    systemNamespace,
	}}

	m.writeValuesJob(vals, &job)

	if err := m.cl.APIClient.Create(ctx, &job); err != nil {
		return err
	}

	if err := m.waitForJobCompletion(ctx, &job); err != nil {
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
	},
		retry.WithBackoffFactory(
			backoff.Build(
				backoff.Exponential(100*time.Millisecond, 2.0),
				backoff.MaxBound(1*time.Minute),
				backoff.MaxRetries(20),
			),
		),
	)
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
