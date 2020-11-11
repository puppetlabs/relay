package dev

import (
	"context"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	rbacmanagerv1beta1 "github.com/fairwindsops/rbac-manager/pkg/apis/rbacmanager/v1beta1"
	certmanagerv1beta1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1beta1"
	certmanagermetav1 "github.com/jetstack/cert-manager/pkg/apis/meta/v1"
	"github.com/puppetlabs/horsehead/v2/workdir"
	installerv1alpha1 "github.com/puppetlabs/relay-core/pkg/apis/install.relay.sh/v1alpha1"
	"github.com/puppetlabs/relay-core/pkg/operator/dependency"
	"github.com/puppetlabs/relay-core/pkg/util/retry"
	v1 "github.com/puppetlabs/relay-core/pkg/workflow/types/v1"
	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/puppetlabs/relay/pkg/dev/manifests"
	"github.com/puppetlabs/relay/pkg/dialog"
	"github.com/puppetlabs/relay/pkg/model"
	helmchartv1 "github.com/rancher/helm-controller/pkg/apis/helm.cattle.io/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage/names"
	kubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	utilflag "k8s.io/component-base/cli/flag"
	kctlcmd "k8s.io/kubernetes/pkg/kubectl/cmd"
	cachingv1alpha1 "knative.dev/caching/pkg/apis/caching/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	DefaultScheme = runtime.NewScheme()
	schemeBuilder = runtime.NewSchemeBuilder(
		kubernetesscheme.AddToScheme,
		metav1.AddMetaToScheme,
		rbacv1.AddToScheme,
		apiextensionsv1.AddToScheme,
		apiextensionsv1beta1.AddToScheme,
		dependency.AddToScheme,
		certmanagerv1beta1.AddToScheme,
		rbacmanagerv1beta1.AddToScheme,
		helmchartv1.AddToScheme,
		cachingv1alpha1.AddToScheme,
		installerv1alpha1.AddToScheme,
	)
	_ = schemeBuilder.AddToScheme(DefaultScheme)
)

const (
	defaultWorkflowName      = "relay-workflow"
	jwtSigningKeysSecretName = "relay-core-v1-operator-signing-keys"
)

type Config struct {
	WorkDir *workdir.WorkDir
	Dialog  dialog.Dialog
}

type Manager struct {
	cm  cluster.Manager
	cl  *cluster.Client
	cfg Config
}

type InitializeOptions struct {
	ImageRegistryPort int
}

func (m *Manager) KubectlCommand() (*cobra.Command, error) {
	if err := os.Setenv("KUBECONFIG", filepath.Join(m.cfg.WorkDir.Path, "kubeconfig")); err != nil {
		return nil, err
	}

	pflag.CommandLine.SetNormalizeFunc(utilflag.WordSepNormalizeFunc)
	pflag.CommandLine.AddGoFlagSet(goflag.CommandLine)

	kubectl := kctlcmd.NewDefaultKubectlCommand()

	return kubectl, nil
}

func (m *Manager) WriteKubeconfig(ctx context.Context) error {
	return m.cm.WriteKubeconfig(ctx, filepath.Join(m.cfg.WorkDir.Path, "kubeconfig"))
}

func (m *Manager) Delete(ctx context.Context) error {
	// TODO fix hack: deletes the PVCs because dirs inside are often created as root
	// and we don't want relay running like that on the host to rm the data dir.
	nm := newNamespaceManager(m.cl)
	if err := nm.delete(ctx, systemNamespace); err != nil && !apierrors.IsNotFound(err) {
		return err
	}

	err := retry.Retry(ctx, 2*time.Second, func() *retry.RetryError {
		pvcs := &corev1.PersistentVolumeClaimList{}
		if err := m.cl.APIClient.List(ctx, pvcs, client.InNamespace(systemNamespace)); err != nil {
			return retry.RetryPermanent(err)
		}

		if len(pvcs.Items) != 0 {
			return retry.RetryTransient(fmt.Errorf("waiting for pvcs to be deleted"))
		}

		return retry.RetryPermanent(nil)
	})
	if err != nil {
		return err
	}

	if err := m.cm.Delete(ctx); err != nil {
		return err
	}

	if err := m.cfg.WorkDir.Cleanup(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) RunWorkflow(ctx context.Context, r io.ReadCloser, params map[string]string) (*model.WorkflowSummary, error) {
	vm := newVaultManager(m.cl, m.cfg)
	am := newAdminManager(m.cl, vm)

	decoder := v1.NewDocumentStreamingDecoder(r, &v1.YAMLDecoder{})

	wd, err := decoder.DecodeStream(ctx)
	if err != nil {
		return nil, err
	}

	name := wd.Name
	if name == "" {
		name = defaultWorkflowName
	}

	runID := names.SimpleNameGenerator.GenerateName(name + "-")
	if err != nil {
		return nil, err
	}

	if err := am.addConnectionForWorkflow(ctx, name); err != nil {
		return nil, err
	}

	runParams := v1.WorkflowRunParameters{}

	for k, v := range params {
		runParams[k] = &v1.WorkflowRunParameter{
			Value: v,
		}
	}

	mapper := v1.NewDefaultRunEngineMapper(
		v1.WithDomainIDRunOption(name),
		v1.WithNamespaceRunOption(name),
		v1.WithWorkflowNameRunOption(name),
		v1.WithWorkflowRunNameRunOption(runID),
		v1.WithVaultEngineMountRunOption("customers"),
		v1.WithRunParametersRunOption(runParams),
	)

	manifest, err := mapper.ToRuntimeObjectsManifest(wd)
	if err != nil {
		return nil, err
	}

	if err := m.cl.APIClient.Create(ctx, manifest.Namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	if err := m.cl.APIClient.Create(ctx, manifest.WorkflowRun); err != nil {
		return nil, err
	}

	ws := &model.WorkflowSummary{
		WorkflowIdentifier: &model.WorkflowIdentifier{
			Name: name,
		},
		Description: wd.Description,
	}

	return ws, nil
}

func (m *Manager) SetWorkflowSecret(ctx context.Context, workflow, key, value string) error {
	vm := newVaultManager(m.cl, m.cfg)
	secret := map[string]string{
		path.Join("customers", "workflows", workflow, key): value,
	}

	return vm.writeSecrets(ctx, secret)
}

func (m *Manager) InitializeRelayCore(ctx context.Context, opts InitializeOptions) error {
	// I introduced some race condition where the cluster hasn't fully setup
	// the object APIs or something, so when we try to create objects here, it
	// will blow up saying the API for that object type doesn't exist. If we
	// sleep for just a second, then we give it enough time to fully warm up or
	// something. I dunno...
	//
	// There's an option in k3d's cluster create that I set to wait for the
	// server, but I think there's something deeper happening inside kubernetes
	// (probably in the API server).
	<-time.After(time.Second)

	log := m.cfg.Dialog

	nm := newNamespaceManager(m.cl)
	vm := newVaultManager(m.cl, m.cfg)
	am := newAdminManager(m.cl, vm)
	rm := newRegistryManager(m.cl)
	rcm := newRelayCoreManager(m.cl)

	if err := nm.create(ctx, systemNamespace); err != nil {
		return err
	}

	if err := nm.create(ctx, registryNamespace); err != nil {
		return err
	}

	log.Info("initializing admin account")
	if err := am.reconcile(ctx); err != nil {
		return err
	}

	log.Info("initializing image registry")
	if err := rm.reconcile(ctx); err != nil {
		return err
	}

	patchers := []objectPatcherFunc{
		nm.objectNamespacePatcher(systemNamespace),
		missingProtocolPatcher,
		registryLoadBalancerPortPatcher(opts.ImageRegistryPort),
	}

	// Apply manifests in ordered phases. Note that some managers
	// have weird dependencies on running services. For instance, you cannot
	// create or apply a ClusterIssuer unless the cert-manager webhook service
	// is Ready. This means we will just wait for all services across all created
	// namespaces to be ready before moving to the next phase of applying manifests.
	// TODO: dynamically generate the list as we process the manifests

	if err := m.processManifests(ctx, "/01-init", patchers, []string{"cert-manager", "relay-system"}); err != nil {
		return err
	}

	if err := rcm.reconcile(ctx); err != nil {
		return err
	}

	if err := m.waitForCertificates(ctx, systemNamespace); err != nil {
		return err
	}

	log.Info("initializing vault")
	if err := vm.reconcileInit(ctx); err != nil {
		return err
	}

	if err := vm.reconcileUnseal(ctx); err != nil {
		return err
	}

	// get the CA secret so we can pass the cert into things that need it.
	caSecretKey := client.ObjectKey{
		Name:      rcm.objects.selfSignedCA.Spec.SecretName,
		Namespace: rcm.objects.selfSignedCA.Namespace,
	}

	tlsSecret := &corev1.Secret{}

	if err := m.cl.APIClient.Get(ctx, caSecretKey, tlsSecret); err != nil {
		return err
	}

	patchers = []objectPatcherFunc{
		nm.objectNamespacePatcher("tekton-pipelines"),
		missingProtocolPatcher,
	}

	if err := m.processManifests(ctx, "/03-tekton", patchers, []string{"tekton-pipelines"}); err != nil {
		return err
	}

	patchers = []objectPatcherFunc{
		nm.objectNamespacePatcher("knative-serving"),
		missingProtocolPatcher,
	}

	if err := m.processManifests(ctx, "/04-knative", patchers, []string{"knative-serving"}); err != nil {
		return err
	}

	patchers = []objectPatcherFunc{
		nm.objectNamespacePatcher(systemNamespace),
		missingProtocolPatcher,
		registryLoadBalancerPortPatcher(opts.ImageRegistryPort),
		admissionPatcher(tlsSecret.Data["ca.crt"]),
	}

	if err := m.processManifests(ctx, "/05-relay", patchers, []string{systemNamespace}); err != nil {
		return err
	}

	if err := vm.reconcileConfiguration(ctx); err != nil {
		return err
	}

	if err := nm.create(ctx, "ambassador-webhook"); err != nil {
		return err
	}

	patchers = []objectPatcherFunc{
		nm.objectNamespacePatcher("ambassador-webhook"),
		missingProtocolPatcher,
		ambassadorPatcher,
	}

	if err := m.processManifests(ctx, "/06-ambassador", patchers, []string{"ambassador-webhook"}); err != nil {
		return err
	}

	patchers = []objectPatcherFunc{
		nm.objectNamespacePatcher("default"),
	}

	if err := m.processManifests(ctx, "/07-hostpath", patchers, nil); err != nil {
		return err
	}

	return nil
}

func (m *Manager) processManifests(ctx context.Context, path string, patchers []objectPatcherFunc, initNamespaces []string) error {
	log := m.cfg.Dialog

	objects, err := m.parseAndLoadManifests(manifests.MustAssetListDir(path)...)
	if err != nil {
		return err
	}

	log.Infof("applying %s", path)
	if err := m.applyAllWithPatchers(ctx, patchers, objects); err != nil {
		return err
	}

	for _, ns := range initNamespaces {
		log.Infof("waiting for services in: %s", ns)
		if err := m.waitForServices(ctx, ns); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) StartRelayCore(ctx context.Context) error {
	log := m.cfg.Dialog
	vm := newVaultManager(m.cl, m.cfg)

	log.Infof("waiting for services in: %s", "relay-system")
	if err := m.waitForServices(ctx, "relay-system"); err != nil {
		return err
	}

	return vm.reconcileUnseal(ctx)
}

func (m *Manager) parseAndLoadManifests(files ...string) ([]runtime.Object, error) {
	log := m.cfg.Dialog
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

func NewManager(cm cluster.Manager, cl *cluster.Client, cfg Config) *Manager {
	return &Manager{
		cm:  cm,
		cl:  cl,
		cfg: cfg,
	}
}
