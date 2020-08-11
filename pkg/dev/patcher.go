package dev

import (
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
