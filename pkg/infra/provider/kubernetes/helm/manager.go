package helm

import (
	"bytes"
	"context"
	"encoding/json"
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
	_, err := m.run(ctx, nil, "init", "--service-account", "tiller", "--tiller-namespace", "tiller-world")

	return err
}

func (m HelmManager) DeployChart(ctx context.Context, release string, chart *Chart, variables map[string]string) errors.Error {
	ok, err := m.releaseExists(ctx, release)
	if err != nil {
		return err
	}

	if ok {
		params := append([]string{release}, chart.params()...)
		_, err = m.run(ctx, variables, "upgrade", params...)
	} else {
		params := append([]string{"-n", release}, chart.params()...)
		_, err = m.run(ctx, variables, "install", params...)
	}

	return err
}

func (m HelmManager) releaseExists(ctx context.Context, release string) (bool, errors.Error) {
	r, err := m.run(ctx, nil, "list", "--output", "json", release)
	if err != nil {
		return false, err
	}

	buf := &bytes.Buffer{}

	buf.ReadFrom(r)

	if buf.Len() == 0 {
		return false, nil
	}

	resp := struct {
		Releases []interface{} `json:"releases"`
	}{}

	if err := json.Unmarshal(buf.Bytes(), &resp); err != nil {
		return false, errors.NewHelmCommandExecError().WithCause(err)
	}

	if len(resp.Releases) == 0 {
		return false, nil
	}

	return true, nil
}

func (m HelmManager) run(ctx context.Context, variables map[string]string, subcmd string, params ...string) (io.Reader, errors.Error) {
	args := append([]string{m.helmcmd, "--kubeconfig", m.kubeconfig, "--tiller-namespace", "tiller-world", subcmd}, params...)
	raw := strings.Join(args, " ")

	out, err := execution.ExecuteCommand(raw, variables, m.logger)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBufferString(out)

	return buf, nil
}

func NewHelmManager(kubeconfig string, logger logging.Logger) *HelmManager {
	return &HelmManager{
		helmcmd:    "helm",
		kubeconfig: kubeconfig,
		logger:     logger,
	}
}
