package dev

import (
	"context"
	"fmt"

	"github.com/puppetlabs/relay/pkg/cluster"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	defaultRegistryPort       = 5000
	defaultRegistryCachePort  = 5001
	defaultRegistryVolumeSize = "5Gi"
	registryVolumeName        = "data"
	registryCacheVolumeName   = "cache-data"
	registryVolumeMountPath   = "/var/lib/registry"
)

type registryObjects struct {
	statefulSet         appsv1.StatefulSet
	headlessService     corev1.Service
	loadBalancerService corev1.Service
}

func newRegistryObjects() *registryObjects {
	objectMeta := metav1.ObjectMeta{
		Name:      "docker-registry",
		Namespace: registryNamespace,
	}

	headlessObjectMeta := objectMeta
	headlessObjectMeta.Name = fmt.Sprintf("%s-headless", headlessObjectMeta.Name)

	return &registryObjects{
		statefulSet:         appsv1.StatefulSet{ObjectMeta: objectMeta},
		headlessService:     corev1.Service{ObjectMeta: headlessObjectMeta},
		loadBalancerService: corev1.Service{ObjectMeta: objectMeta},
	}
}

type registryManager struct {
	cl         *cluster.Client
	objects    *registryObjects
	serverPort int32
	cachePort  int32
}

func (m *registryManager) reconcile(ctx context.Context) error {
	client := m.cl.APIClient

	if _, err := ctrl.CreateOrUpdate(ctx, client, &m.objects.statefulSet, func() error {
		m.statefulSet(&m.objects.statefulSet)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, client, &m.objects.headlessService, func() error {
		m.headlessService(&m.objects.headlessService)

		return nil
	}); err != nil {
		return err
	}

	if _, err := ctrl.CreateOrUpdate(ctx, client, &m.objects.loadBalancerService, func() error {
		m.loadBalancerService(&m.objects.loadBalancerService)

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (m *registryManager) statefulSet(statefulSet *appsv1.StatefulSet) {
	labels := m.labels()
	replicas := int32(1)

	statefulSet.Labels = labels
	statefulSet.Spec.Selector = &metav1.LabelSelector{MatchLabels: labels}
	statefulSet.Spec.Template.Labels = labels

	statefulSet.Spec.ServiceName = m.objects.headlessService.Name
	statefulSet.Spec.Replicas = &replicas

	template := &statefulSet.Spec.Template.Spec

	cacheContainer := corev1.Container{
		Name:  "cache",
		Image: "registry:2",
		Ports: []corev1.ContainerPort{
			{Name: "cache-http", ContainerPort: m.cachePort, Protocol: corev1.ProtocolTCP},
		},
		VolumeMounts: []corev1.VolumeMount{
			{Name: registryCacheVolumeName, MountPath: registryVolumeMountPath},
		},
	}

	m.cacheContainerEnv(&cacheContainer)

	template.Containers = []corev1.Container{
		{
			Name:  "server",
			Image: "registry:2",
			Ports: []corev1.ContainerPort{
				{Name: "http", ContainerPort: m.serverPort, Protocol: corev1.ProtocolTCP},
			},
			VolumeMounts: []corev1.VolumeMount{
				{Name: registryVolumeName, MountPath: registryVolumeMountPath},
			},
		},
		cacheContainer,
	}

	dataClaim := corev1.PersistentVolumeClaim{
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{corev1.ReadWriteOnce},
			Resources: corev1.ResourceRequirements{
				Requests: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceStorage: resource.MustParse(defaultRegistryVolumeSize),
				},
			},
		},
	}
	dataClaim.Name = registryVolumeName

	cacheClaim := corev1.PersistentVolumeClaim{Spec: dataClaim.Spec}
	cacheClaim.Name = registryCacheVolumeName

	statefulSet.Spec.VolumeClaimTemplates = []corev1.PersistentVolumeClaim{dataClaim, cacheClaim}
}

func (m *registryManager) cacheContainerEnv(container *corev1.Container) {
	container.Env = append(container.Env,
		corev1.EnvVar{
			Name:  "REGISTRY_PROXY_REMOTEURL",
			Value: "https://registry-1.docker.io",
		},
		corev1.EnvVar{
			Name:  "REGISTRY_HTTP_ADDR",
			Value: fmt.Sprintf("0.0.0.0:%d", m.cachePort),
		},
		corev1.EnvVar{
			Name:  "REGISTRY_LOG_LEVEL",
			Value: "debug",
		},
	)
}

func (m *registryManager) headlessService(service *corev1.Service) {
	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       m.serverPort,
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "cache-http",
			Port:       m.cachePort,
			TargetPort: intstr.FromString("cache-http"),
			Protocol:   corev1.ProtocolTCP,
		},
	}

	service.Spec.ClusterIP = "None"
	service.Spec.Selector = m.labels()
}

func (m *registryManager) loadBalancerService(service *corev1.Service) {
	service.Spec.Ports = []corev1.ServicePort{
		{
			Name:       "http",
			Port:       m.serverPort,
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "cache-http",
			Port:       m.cachePort,
			TargetPort: intstr.FromString("cache-http"),
			Protocol:   corev1.ProtocolTCP,
		},
	}

	service.Spec.Type = corev1.ServiceTypeLoadBalancer
	service.Spec.Selector = m.labels()
}

func (m *registryManager) labels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/name":      "docker-registry",
		"app.kubernetes.io/component": "server",
	}
}

func newRegistryManager(cl *cluster.Client) *registryManager {
	return &registryManager{
		cl:         cl,
		objects:    newRegistryObjects(),
		serverPort: defaultRegistryPort,
		cachePort:  defaultRegistryCachePort,
	}
}
