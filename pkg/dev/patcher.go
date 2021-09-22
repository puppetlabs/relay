package dev

import (
	"github.com/puppetlabs/leg/k8sutil/pkg/manifest"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func registryLoadBalancerPortPatcher(registryPort int) manifest.PatcherFunc {
	return func(obj manifest.Object, gvk *schema.GroupVersionKind) {
		service, ok := obj.(*corev1.Service)
		if ok && service.Name == "docker-registry" && service.Namespace == "docker-registry" {
			for i, port := range service.Spec.Ports {
				if port.Name == "http" {
					service.Spec.Ports[i].Port = int32(registryPort)

					break
				}
			}
		}
	}
}

func ambassadorPatcher() manifest.PatcherFunc {
	return func(obj manifest.Object, gvk *schema.GroupVersionKind) {
		deployment, ok := obj.(*appsv1.Deployment)
		if !ok || deployment.GetName() != "ambassador" {
			return
		}

		for i, c := range deployment.Spec.Template.Spec.Containers {
			if c.Name != "ambassador" {
				continue
			}

			setKubernetesEnvVar(&c.Env, "AMBASSADOR_ID", "webhook")
			setKubernetesEnvVar(&c.Env, "AMBASSADOR_KNATIVE_SUPPORT", "true")

			deployment.Spec.Template.Spec.Containers[i] = c
		}

		// Make as minimal as possible for testing.
		deployment.Spec.Replicas = func(i int32) *int32 { return &i }(1)
		deployment.Spec.RevisionHistoryLimit = func(i int32) *int32 { return &i }(0)

		// Don't allow old pods to linger.
		deployment.Spec.Strategy.Type = appsv1.RecreateDeploymentStrategyType
		deployment.Spec.Strategy.RollingUpdate = nil
	}
}

func setKubernetesEnvVar(target *[]corev1.EnvVar, name, value string) {
	for i, ev := range *target {
		if ev.Name == name {
			(*target)[i].Value = value
			return
		}
	}

	*target = append(*target, corev1.EnvVar{Name: name, Value: value})
}
