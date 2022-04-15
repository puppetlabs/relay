package dev

import (
	"context"
	"io"
	"path"
	"time"

	"github.com/puppetlabs/leg/workdir"
	v1 "github.com/puppetlabs/relay-client-go/models/pkg/workflow/types/v1"
	installerv1alpha1 "github.com/puppetlabs/relay-core/pkg/apis/install.relay.sh/v1alpha1"
	relayv1beta1 "github.com/puppetlabs/relay-core/pkg/apis/relay.sh/v1beta1"
	"github.com/puppetlabs/relay-core/pkg/obj"
	"github.com/puppetlabs/relay-core/pkg/operator/dependency"
	helmchartv1 "github.com/rancher/helm-controller/pkg/apis/helm.cattle.io/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage/names"
	kubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	cachingv1alpha1 "knative.dev/caching/pkg/apis/caching/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
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
		helmchartv1.AddToScheme,
		cachingv1alpha1.AddToScheme,
		installerv1alpha1.AddToScheme,
	)
	_ = schemeBuilder.AddToScheme(DefaultScheme)
)

const (
	defaultWorkflowName      = "relay-workflow"
	jwtSigningKeysSecretName = "relay-core-v1-operator-signing-keys"

	VaultEngineMountCustomers = "customers"
	VaultEngineMountWorkflows = "workflows"
)

type Client struct {
	APIClient client.Client
	Mapper    meta.RESTMapper
}

type ClientOptions struct {
	Scheme *runtime.Scheme
}

type Config struct {
	WorkDir *workdir.WorkDir
}

type Manager struct {
	cl  *Client
	cfg Config
}

type InitializeOptions struct {
	InstallHelmController bool
}

type InstallerOptions struct {
	InstallerImage                            string
	LogServiceImage                           string
	MetadataAPIImage                          string
	OperatorImage                             string
	OperatorVaultInitImage                    string
	OperatorWebhookCertificateControllerImage string
	VaultServerImage                          string
	VaultSidecarImage                         string
}

// FIXME Consider a better mechanism for specific service options
type LogServiceOptions struct {
	Enabled               bool
	CredentialsKey        string
	CredentialsSecretName string
	Project               string
	Dataset               string
	Table                 string
}

func (m *Manager) LoadWorkflow(ctx context.Context, r io.ReadCloser) (*v1.WorkflowData, error) {
	decoder := v1.NewDocumentStreamingDecoder(r, &v1.YAMLDecoder{})

	wd, err := decoder.DecodeStream(ctx)
	if err != nil {
		return nil, err
	}

	return wd, nil
}

func (m *Manager) CreateTenant(ctx context.Context, name string) (*relayv1beta1.Tenant, error) {
	mapper := v1.NewDefaultTenantEngineMapper(
		v1.WithNameTenantOption(name),
		v1.WithIDTenantOption(name),
		v1.WithWorkflowNameTenantOption(name),
		v1.WithWorkflowIDTenantOption(name),
		v1.WithNamespaceTenantOption(tenantNamespace),
	)

	mapping, err := mapper.ToRuntimeObjectsManifest()
	if err != nil {
		return nil, err
	}

	if err := m.cl.APIClient.Create(ctx, mapping.Namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	key := client.ObjectKey{
		Name:      mapping.Tenant.GetName(),
		Namespace: mapping.Tenant.GetNamespace(),
	}

	t := obj.NewTenant(key)
	if _, err := t.Load(ctx, m.cl.APIClient); err != nil {
		return nil, err
	}

	t.Object.Spec = mapping.Tenant.Spec

	if err := t.Persist(ctx, m.cl.APIClient); err != nil {
		return nil, err
	}

	return mapping.Tenant, nil
}

func (m *Manager) CreateWorkflow(ctx context.Context, wd *v1.WorkflowData, t *relayv1beta1.Tenant) (*relayv1beta1.Workflow, error) {
	vm := newVaultManager(m.cl, m.cfg)
	am := newAdminManager(m.cl, vm)

	name := wd.Name
	if name == "" {
		name = defaultWorkflowName
	}

	// FIXME Refactor the connection handling (ideally not directly linked to the create workflow functionality)
	if err := am.reconcile(ctx); err != nil {
		return nil, err
	}
	if err := am.addConnectionForWorkflow(ctx, name); err != nil {
		return nil, err
	}

	mapper := v1.NewDefaultWorkflowMapper(
		v1.WithDomainIDOption(name),
		v1.WithNamespaceOption(tenantNamespace),
		v1.WithWorkflowNameOption(name),
		v1.WithVaultEngineMountOption(VaultEngineMountCustomers),
		v1.WithTenantOption(t),
	)

	mapping, err := mapper.Map(wd)
	if err != nil {
		return nil, err
	}

	key := client.ObjectKey{
		Name:      mapping.Workflow.GetName(),
		Namespace: mapping.Workflow.GetNamespace(),
	}

	wf := obj.NewWorkflow(key)
	if _, err := wf.Load(ctx, m.cl.APIClient); err != nil {
		return nil, err
	}

	wf.Object.Spec = mapping.Workflow.Spec

	if err := wf.Persist(ctx, m.cl.APIClient); err != nil {
		return nil, err
	}

	return mapping.Workflow, nil
}

func (m *Manager) RunWorkflow(ctx context.Context, wf *relayv1beta1.Workflow, params map[string]string) (*relayv1beta1.Run, error) {
	runName := names.SimpleNameGenerator.GenerateName(wf.GetName() + "-")

	runParams := v1.WorkflowRunParameters{}

	for k, v := range params {
		runParams[k] = &v1.WorkflowRunParameter{
			Value: v,
		}
	}

	mapper := v1.NewDefaultRunEngineMapper(
		v1.WithDomainIDRunOption(wf.GetName()),
		v1.WithNamespaceRunOption(wf.GetNamespace()),
		v1.WithWorkflowNameRunOption(wf.GetName()),
		v1.WithWorkflowRunNameRunOption(runName),
		v1.WithVaultEngineMountRunOption(VaultEngineMountCustomers),
		v1.WithRunParametersRunOption(runParams),
		v1.WithWorkflowRunOption(wf),
	)

	mapping, err := mapper.ToRuntimeObjectsManifest()
	if err != nil {
		return nil, err
	}

	if err := m.cl.APIClient.Create(ctx, mapping.Namespace); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	if err := m.cl.APIClient.Create(ctx, mapping.WorkflowRun); err != nil {
		return nil, err
	}

	return mapping.WorkflowRun, err
}

func (m *Manager) SetWorkflowSecret(ctx context.Context, workflow, key, value string) error {
	vm := newVaultManager(m.cl, m.cfg)
	secret := map[string]string{
		path.Join(VaultEngineMountCustomers, VaultEngineMountWorkflows, workflow, key): value,
	}

	return vm.writeSecrets(ctx, secret)
}

func (m *Manager) InitializeRelayCore(ctx context.Context, initOpts InitializeOptions, installerOpts InstallerOptions, logServiceOpts LogServiceOptions) error {
	// I introduced some race condition where the cluster hasn't fully setup
	// the object APIs or something, so when we try to create objects here, it
	// will blow up saying the API for that object type doesn't exist. If we
	// sleep for just a second, then we give it enough time to fully warm up or
	// something. I dunno...
	//
	// There's an option in k3d's cluster create that I set to wait for the
	// server, but I think there's something deeper happening inside kubernetes
	// (probably in the API server).
	<-time.After(time.Second * 5)

	nm := newNamespaceManager(m.cl)
	rim := newRelayInstallerManager(m.cl, installerOpts)
	rcm := newRelayCoreManager(m.cl, installerOpts, logServiceOpts)

	if err := nm.reconcile(ctx); err != nil {
		return err
	}

	// TODO: dynamically generate the list as we process the manifests

	mm := NewManifestManager(m.cl)

	manifests := []string{
		"/tekton",
		"/knative",
		"/relay",
		"/kourier",
	}

	if initOpts.InstallHelmController {
		manifests = append(manifests, "helm-controller")
	}

	for _, manifest := range manifests {
		if err := mm.ProcessManifests(ctx, manifest); err != nil {
			return err
		}
	}

	if err := rim.reconcile(ctx); err != nil {
		return err
	}

	if err := rcm.reconcile(ctx); err != nil {
		return err
	}

	return nil
}

func NewManager(ctx context.Context) (*Manager, error) {
	kcfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	apiConfig, err := kcfg.RawConfig()
	if err != nil {
		return nil, err
	}

	cl, err := NewClient(ctx, ClientOptions{Scheme: DefaultScheme}, &apiConfig)
	if err != nil {
		return nil, err
	}

	return &Manager{
		cl:  cl,
		cfg: Config{},
	}, nil
}

func NewClient(ctx context.Context, opts ClientOptions, apiConfig *clientcmdapi.Config) (*Client, error) {
	overrides := &clientcmd.ConfigOverrides{}
	clientConfig := clientcmd.NewDefaultClientConfig(*apiConfig, overrides)

	restConfig, err := clientConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	c, err := client.New(restConfig, client.Options{
		Scheme: opts.Scheme,
	})
	if err != nil {
		return nil, err
	}

	mapper, err := apiutil.NewDynamicRESTMapper(restConfig)
	if err != nil {
		return nil, err
	}

	return &Client{
		APIClient: c,
		Mapper:    mapper,
	}, nil
}
