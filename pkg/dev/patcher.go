package dev

import (
	admissionregistrationv1beta1 "k8s.io/api/admissionregistration/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type objectPatcherFunc func(runtime.Object)

// TODO replace when relay-core k8s testutils are generalized
func missingProtocolPatcher(obj runtime.Object) {
	switch t := obj.(type) {
	case *appsv1.Deployment:
		// SSA has marked "protocol" is required but basically everyone expects
		// it to default to TCP.
		for i, container := range t.Spec.Template.Spec.Containers {
			for j, port := range container.Ports {
				if len(port.Protocol) > 0 {
					continue
				}

				t.Spec.Template.Spec.Containers[i].Ports[j].Protocol = corev1.ProtocolTCP
			}
		}
	case *corev1.Service:
		// Same for services.
		for i, port := range t.Spec.Ports {
			if len(port.Protocol) > 0 {
				continue
			}

			t.Spec.Ports[i].Protocol = corev1.ProtocolTCP
		}
	}
}

func registryLoadBalancerPortPatcher(registryPort int) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *corev1.Service:
			if t.Name == "docker-registry" && t.Namespace == "docker-registry" {
				for i, port := range t.Spec.Ports {
					if port.Name == "http" {
						t.Spec.Ports[i].Port = int32(registryPort)

						break
					}
				}
			}
		}
	}
}

func ambassadorPatcher(obj runtime.Object) {
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

func admissionPatcher(caCertData []byte) objectPatcherFunc {
	return func(obj runtime.Object) {
		switch t := obj.(type) {
		case *admissionregistrationv1beta1.MutatingWebhookConfiguration:
			for i := range t.Webhooks {
				t.Webhooks[i].ClientConfig.CABundle = caCertData
			}
		}
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
