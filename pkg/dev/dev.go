package dev

import (
	"context"
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	rbacmanagerv1beta1 "github.com/fairwindsops/rbac-manager/pkg/apis/rbacmanager/v1beta1"
	certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
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
	scheme        = runtime.NewScheme()
	schemeBuilder = runtime.NewSchemeBuilder(
		kubernetesscheme.AddToScheme,
		metav1.AddMetaToScheme,
		apiextensionsv1beta1.AddToScheme,
		dependency.AddToScheme,
		certmanagerv1beta1.AddToScheme,
		rbacmanagerv1beta1.AddToScheme,
		helmchartv1.AddToScheme,
	)
	_ = schemeBuilder.AddToScheme(scheme)
)

type Options struct {
	DataDir string
}

type Manager struct {
	cm   cluster.Manager
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

func (m *Manager) ApplyCoreResources(ctx context.Context) error {
	cl, err := m.cm.GetClient(ctx, cluster.ClientOptions{
		Scheme: scheme,
	})
	if err != nil {
		return err
	}

	nm := newNamespaceManager(cl)
	cam := newCAManager()

	if err := nm.create(ctx); err != nil {
		return err
	}

	caPair, err := cam.generateCA()
	if err != nil {
		return err
	}

	patchers := []objectPatcherFunc{
		nm.objectNamespacePatcher("system"),
		missingProtocolPatcher,
		cam.secretPatcher("relay-cert-manager-ca", caPair),
		cam.certificatePatcher("relay-cluster-issuer"),
		cam.admissionPatcher(caPair),
	}

	// Manifests are split into diffent directories because some managers
	// have weird dependencies on running services. For instance, you cannot
	// create or apply a ClusterIssuer unless the cert-manager webhook service
	// is Ready. This means we will just wait for all services across all created
	// namespaces to be ready before moving to the next phase of applying manifests.
	initManifests := manifests.MustAssetListDir("/01-init")
	initObjects := []runtime.Object{}
	// TODO: dynamically generate this list as we proccess the manifests
	initNamespaces := []string{"cert-manager", "tekton-pipelines"}

	for _, f := range initManifests {
		manifest := manifests.MustAsset(f)

		log.Infof("parsing manifest %s", f)

		objs, err := parseManifest(manifest)
		if err != nil {
			return err
		}

		initObjects = append(initObjects, objs...)
	}

	log.Info("applying init objects")
	if err := m.applyAllWithPatchers(ctx, cl, patchers, initObjects); err != nil {
		return err
	}

	for _, ns := range initNamespaces {
		log.Infof("waiting for services in: %s", ns)
		if err := m.waitForServices(ctx, cl, ns); err != nil {
			return err
		}
	}

	relayManifests := manifests.MustAssetListDir("/03-relay")
	relayObjects := []runtime.Object{}

	for _, f := range relayManifests {
		manifest := manifests.MustAsset(f)

		log.Infof("parsing manifest %s", f)

		objs, err := parseManifest(manifest)
		if err != nil {
			return err
		}

		relayObjects = append(relayObjects, objs...)
	}

	log.Info("applying relay objects")
	if err := m.applyAllWithPatchers(ctx, cl, patchers, relayObjects); err != nil {
		return err
	}

	return nil
}

func (m *Manager) waitForServices(ctx context.Context, cl *cluster.Client, namespace string) error {
	err := retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
		eps := &corev1.EndpointsList{}
		if err := cl.APIClient.List(ctx, eps, client.InNamespace(namespace)); err != nil {
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

func (m *Manager) apply(ctx context.Context, cl *cluster.Client, obj runtime.Object) error {
	if err := cl.APIClient.Patch(ctx, obj, client.Apply, client.ForceOwnership, client.FieldOwner("relay-e2e")); err != nil {
		return fmt.Errorf("failed to apply object '%s': %w", obj.GetObjectKind().GroupVersionKind().String(), err)
	}

	return nil
}

func (m *Manager) applyAllWithPatchers(ctx context.Context, cl *cluster.Client, patchers []objectPatcherFunc, objs []runtime.Object) error {
	for _, obj := range objs {
		for _, patcher := range patchers {
			patcher(obj)
		}

		if err := m.apply(ctx, cl, obj); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) DeleteDataDir() error {
	return os.RemoveAll(m.opts.DataDir)
}

func (m *Manager) kubectlExec(args ...string) error {
	kubectl, err := m.KubectlCommand()
	if err != nil {
		return err
	}

	kubectl.SetArgs(args)

	return kubectl.Execute()
}

func NewManager(cm cluster.Manager, opts Options) *Manager {
	return &Manager{
		cm:   cm,
		opts: opts,
	}
}
