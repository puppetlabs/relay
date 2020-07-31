package dev

import (
	"context"
	goflag "flag"
	"fmt"
	"os"
	"path/filepath"

	certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	"github.com/puppetlabs/relay-core/pkg/dependency"
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev/manifests"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
		dependency.AddToScheme,
		certmanagerv1beta1.AddToScheme,
		apiextensionsv1beta1.AddToScheme,
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
	cm := newCAManager(cl, nm.objectNamespacePatcher("system"))

	if err := nm.create(ctx); err != nil {
		return err
	}

	if err := cm.create(ctx); err != nil {
		return err
	}

	files, err := manifests.AssetListDir()
	if err != nil {
		return err
	}

	manifest := manifests.MustAsset(files[0])

	objs, err := parseManifest(manifest)
	if err != nil {
		return err
	}

	systemPatcher := nm.objectNamespacePatcher("system")

	for _, obj := range objs {
		systemPatcher(obj)
		if err := m.apply(ctx, cl, obj); err != nil {
			return fmt.Errorf("failed to apply object '%s': %w", obj.GetObjectKind().GroupVersionKind().String(), err)
		}
	}

	return nil
}

func (m *Manager) apply(ctx context.Context, cl *cluster.Client, obj runtime.Object) error {
	return cl.APIClient.Patch(ctx, obj, client.Apply, client.FieldOwner("relay-e2e"))
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

type objectPatcherFunc func(runtime.Object)
