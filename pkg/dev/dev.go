package dev

import (
	"context"
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	rbacmanagerv1beta1 "github.com/fairwindsops/rbac-manager/pkg/apis/rbacmanager/v1beta1"
	certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	certmanagermetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/puppetlabs/relay-core/pkg/dependency"
	"github.com/puppetlabs/relay-core/pkg/util/retry"
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev/manifests"
	helmchartv1 "github.com/rancher/helm-controller/pkg/apis/helm.cattle.io/v1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	kubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	utilflag "k8s.io/component-base/cli/flag"
	kctlcmd "k8s.io/kubernetes/pkg/kubectl/cmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	DefaultScheme = runtime.NewScheme()
	schemeBuilder = runtime.NewSchemeBuilder(
		kubernetesscheme.AddToScheme,
		metav1.AddMetaToScheme,
		apiextensionsv1beta1.AddToScheme,
		dependency.AddToScheme,
		certmanagerv1beta1.AddToScheme,
		rbacmanagerv1beta1.AddToScheme,
		helmchartv1.AddToScheme,
	)
	_          = schemeBuilder.AddToScheme(DefaultScheme)
	coreImages = []string{
		"relaysh/relay-operator:latest",
		"relaysh/relay-metadata-api:latest",
	}
)

type Options struct {
	DataDir string
}

type Manager struct {
	cm   cluster.Manager
	cl   *cluster.Client
	opts Options
}

func (m *Manager) KubectlCommand() (*cobra.Command, error) {
	if err := os.Setenv("KUBECONFIG", filepath.Join(m.opts.DataDir, "kubeconfig")); err != nil {
		return nil, err
	}

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	kubectl := kctlcmd.NewDefaultKubectlCommand()

	return kubectl, nil
}

func (m *Manager) WriteKubeconfig(ctx context.Context) error {
	if err := os.MkdirAll(m.opts.DataDir, 0700); err != nil {
		return err
	}

	return m.cm.WriteKubeconfig(ctx, filepath.Join(m.opts.DataDir, "kubeconfig"))
}

func (m *Manager) DeleteDataDir() error {
	return os.RemoveAll(m.opts.DataDir)
}

func (m *Manager) InitializeRelayCore(ctx context.Context) error {
	nm := newNamespaceManager(m.cl)
	cam := newCAManager(m.cl)

	// attempting to import host images in the container runtime allows us to
	// quickly bootstrap cluster with custom relay-core builds. if no such
	// images exist, then they will be pulled from the remote by the cluster.
	log.Info("importing locally cached core image")
	if err := m.importInitialImages(ctx); err != nil {
		return err
	}

	if err := nm.create(ctx); err != nil {
		return err
	}

	patchers := []objectPatcherFunc{
		nm.objectNamespacePatcher("system"),
		missingProtocolPatcher,
	}

	// Manifests are split into diffent directories because some managers
	// have weird dependencies on running services. For instance, you cannot
	// create or apply a ClusterIssuer unless the cert-manager webhook service
	// is Ready. This means we will just wait for all services across all created
	// namespaces to be ready before moving to the next phase of applying manifests.
	initObjects, err := m.parseAndLoadManifests(manifests.MustAssetListDir("/01-init")...)
	if err != nil {
		return err
	}

	log.Info("applying init objects")
	if err := m.applyAllWithPatchers(ctx, patchers, initObjects); err != nil {
		return err
	}

	// TODO: dynamically generate this list as we proccess the manifests
	initNamespaces := []string{"cert-manager", "tekton-pipelines"}

	for _, ns := range initNamespaces {
		log.Infof("waiting for services in: %s", ns)
		if err := m.waitForServices(ctx, ns); err != nil {
			return err
		}
	}

	secretObjects, err := m.parseAndLoadManifests(manifests.MustAssetListDir("/02-secrets")...)
	if err != nil {
		return err
	}

	log.Info("applying secret objects")
	if err := m.applyAllWithPatchers(ctx, patchers, secretObjects); err != nil {
		return err
	}

	if err := m.waitForCertificates(ctx, nm.getByID("system")); err != nil {
		return err
	}

	// get the CA secret so we can pass the cert into things that need it.
	caSecretKey := client.ObjectKey{
		Name:      "relay-cert-ca-tls",
		Namespace: nm.getByID("system"),
	}

	tlsSecret := &corev1.Secret{}

	if err := m.cl.APIClient.Get(ctx, caSecretKey, tlsSecret); err != nil {
		return err
	}

	patchers = append(patchers, cam.admissionPatcher(tlsSecret.Data["ca.crt"]))

	relayObjects, err := m.parseAndLoadManifests(manifests.MustAssetListDir("/03-relay")...)
	if err != nil {
		return err
	}

	log.Info("applying relay objects")
	if err := m.applyAllWithPatchers(ctx, patchers, relayObjects); err != nil {
		return err
	}

	return nil
}

func (m *Manager) importInitialImages(ctx context.Context) error {
	for _, image := range coreImages {
		if err := m.cm.ImportImage(ctx, image); err != nil {
			// ignores not found errors
			if strings.Contains(err.Error(), "No valid images specified") {
				continue
			}

			return err
		}
	}

	return nil
}

func (m *Manager) parseAndLoadManifests(files ...string) ([]runtime.Object, error) {
	objects := []runtime.Object{}

	for _, f := range files {
		manifest := manifests.MustAsset(f)

		log.Infof("parsing manifest %s", f)

		manifestObjects, err := parseManifest(manifest)
		if err != nil {
			return nil, err
		}

		objects = append(objects, manifestObjects...)
	}

	return objects, nil
}

func (m *Manager) waitForServices(ctx context.Context, namespace string) error {
	err := retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
		eps := &corev1.EndpointsList{}
		if err := m.cl.APIClient.List(ctx, eps, client.InNamespace(namespace)); err != nil {
			return retry.RetryPermanent(err)
		}

		if len(eps.Items) == 0 {
			return retry.RetryTransient(fmt.Errorf("waiting for endpoints"))
		}

		for _, ep := range eps.Items {
			log.Infof("checking service %s", ep.Name)
			if len(ep.Subsets) == 0 {
				return retry.RetryTransient(fmt.Errorf("waiting for subsets"))
			}

			for _, subset := range ep.Subsets {
				if len(subset.Addresses) == 0 {
					return retry.RetryTransient(fmt.Errorf("waiting for pod assignment"))
				}
			}
		}

		return retry.RetryPermanent(nil)
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) waitForCertificates(ctx context.Context, namespace string) error {
	err := retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
		certs := &certmanagerv1beta1.CertificateList{}
		if err := m.cl.APIClient.List(ctx, certs, client.InNamespace(namespace)); err != nil {
			return retry.RetryPermanent(err)
		}

		if len(certs.Items) == 0 {
			return retry.RetryTransient(fmt.Errorf("waiting for certificates"))
		}

		for _, cert := range certs.Items {
			var isReady bool

			for _, cond := range cert.Status.Conditions {
				if cond.Type == certmanagerv1beta1.CertificateConditionReady {
					isReady = cond.Status == certmanagermetav1.ConditionTrue
				}
			}

			if !isReady {
				return retry.RetryTransient(fmt.Errorf("waiting for certificates to be ready"))
			}
		}

		return retry.RetryPermanent(nil)
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) apply(ctx context.Context, obj runtime.Object) error {
	if err := m.cl.APIClient.Patch(ctx, obj, client.Apply, client.ForceOwnership, client.FieldOwner("relay-e2e")); err != nil {
		return fmt.Errorf("failed to apply object '%s': %w", obj.GetObjectKind().GroupVersionKind().String(), err)
	}

	return nil
}

func (m *Manager) applyAllWithPatchers(ctx context.Context, patchers []objectPatcherFunc, objs []runtime.Object) error {
	for _, obj := range objs {
		for _, patcher := range patchers {
			patcher(obj)
		}

		if err := m.apply(ctx, obj); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) kubectlExec(args ...string) error {
	kubectl, err := m.KubectlCommand()
	if err != nil {
		return err
	}

	kubectl.SetArgs(args)

	return kubectl.Execute()
}

func NewManager(cm cluster.Manager, cl *cluster.Client, opts Options) *Manager {
	return &Manager{
		cm:   cm,
		cl:   cl,
		opts: opts,
	}
}
