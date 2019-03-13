package helm

import (
	"context"
	"io"
	"strings"

	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/execution"
)

type HelmManager struct {
	helmcmd    string
	kubeconfig string
	logger     logging.Logger
}

func (m HelmManager) InitTiller(ctx context.Context) errors.Error {
	_, err := m.run(ctx, "init")

	return err
}

func (m HelmManager) DeployChart(ctx context.Context, release string, chart *Chart) errors.Error {
	params := append([]string{release}, chart.params()...)
	_, err := m.run(ctx, "upgrade", params...)

	return err
}

func (m HelmManager) run(ctx context.Context, subcmd string, params ...string) (io.Reader, errors.Error) {
	args := append([]string{m.helmcmd, "--kubeconfig", m.kubeconfig, subcmd}, params...)
	raw := strings.Join(args, " ")

	if err := execution.ExecuteCommand(raw, nil, m.logger); err != nil {
		return nil, err
	}

	// cmd := exec.CommandContext(ctx, m.helmcmd, args...)

	// if err := cmd.Run(); err != nil {
	// 	return nil, errors.NewHelmCommandExecError().WithCause(err).Bug()
	// }

	return nil, nil
}

func NewHelmManager(kubeconfig string, logger logging.Logger) *HelmManager {
	return &HelmManager{
		helmcmd:    "helm",
		kubeconfig: kubeconfig,
		logger:     logger,
	}
}
