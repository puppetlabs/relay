package dev

import (
	"context"
	goflag "flag"
	"os"
	"path/filepath"

	"github.com/puppetlabs/relay/pkg/cluster"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	utilflag "k8s.io/component-base/cli/flag"
	kctlcmd "k8s.io/kubernetes/pkg/kubectl/cmd"
)

const (
	tektonResourceURL    = "https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.12.0/release.yaml"
	relayCoreResourceURL = "https://raw.githubusercontent.com/puppetlabs/relay-core/master/manifests/resources/nebula.puppet.com_workflowruns.yaml"
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
	if err := m.kubectlExec("apply", "-f", tektonResourceURL); err != nil {
		return err
	}

	return m.kubectlExec("apply", "-f", relayCoreResourceURL)
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
