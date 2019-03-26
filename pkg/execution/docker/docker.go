package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/google/uuid"
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/execution"
	"github.com/puppetlabs/nebula/pkg/plan/encoding"
	plantypes "github.com/puppetlabs/nebula/pkg/plan/types"
)

const (
	defaultDockerHostSocketPath   = "/var/run/docker.sock"
	defaultRunnerMountHostPrefix  = "nebula-runner-context-"
	defaultRunnerContextMountPath = "/nebula-runner-context"
	defaultBuildMountPath         = "/build"
)

type containerConfig struct {
	name      string
	container *container.Config
	host      *container.HostConfig
	network   *network.NetworkingConfig
}

type ExecutorOptions struct {
	DockerHostSocketPath string
	Registry             string
	RegistryUser         string
	RegistryPass         string
}

type Executor struct {
	client       *client.Client
	socketPath   string
	registryAuth types.AuthConfig
}

func (e *Executor) Kind() string {
	return "docker"
}

func (e *Executor) ScheduleAction(ctx context.Context, r execution.ExecutorRuntime,
	action *plantypes.Action, env map[string]string) errors.Error {

	if err := e.pullImage(ctx, r, action.Image); err != nil {
		return err
	}

	config, err := e.createConfig(action, env)
	if err != nil {
		return err
	}

	config.name = e.actionRunName(action)

	id, err := e.createContainer(ctx, r, config)
	if err != nil {
		return err
	}

	if err := e.client.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
		return errors.NewDockerContainerStartError().WithCause(err)
	}

	if err := e.setupLogging(ctx, r, id); err != nil {
		return err
	}

	statusCh, errCh := e.client.ContainerWait(ctx, id, container.WaitConditionNotRunning)

	select {
	case err := <-errCh:
		return errors.NewDockerContainerExecutionError().WithCause(err)
	case st := <-statusCh:
		if st.StatusCode > 0 {
			return errors.NewDockerContainerExecutionError()
		}
	}

	return nil
}

func (e *Executor) actionRunName(action *plantypes.Action) string {
	uid := uuid.New()

	return fmt.Sprintf("nebula-action-%s-%s", action.Name, uid.String())
}

func (e *Executor) pullImage(ctx context.Context, r execution.ExecutorRuntime, image string) errors.Error {
	var authStr string

	if e.registryAuth.Username != "" && e.registryAuth.Password != "" {
		b, err := json.Marshal(e.registryAuth)
		if err != nil {
			return errors.NewDockerCredentialEncodingError().WithCause(err)
		}

		authStr = base64.URLEncoding.EncodeToString(b)
	}

	reader, err := e.client.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: authStr})
	if err != nil {
		return errors.NewDockerImagePullError(image).WithCause(err)
	}

	if _, err := io.Copy(r.IO().Out, reader); err != nil {
		return errors.NewDockerImagePullError(image).WithCause(err)
	}

	if err := reader.Close(); err != nil {
		return errors.NewDockerImagePullError(image).WithCause(err)
	}

	r.Logger().Debug("action-image-pulled", "image", image)

	return nil
}

func (e *Executor) setupRunnerMount() (string, errors.Error) {
	tmpdir, err := ioutil.TempDir("", defaultRunnerMountHostPrefix)
	if err != nil {
		return "", errors.NewDockerUnknownError().WithCause(err).Bug()
	}

	return fmt.Sprintf("%s:%s", tmpdir, defaultRunnerContextMountPath), nil
}

func (e *Executor) setupBuildMount() (string, errors.Error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", errors.NewDockerUnknownError().WithCause(err)
	}

	return fmt.Sprintf("%s:%s", pwd, defaultBuildMountPath), nil
}

func (e *Executor) createConfig(action *plantypes.Action, env map[string]string) (containerConfig, errors.Error) {
	runnerMount, err := e.setupRunnerMount()
	if err != nil {
		return containerConfig{}, err
	}

	buildMount, err := e.setupBuildMount()
	if err != nil {
		return containerConfig{}, err
	}

	encodedSpec, err := encoding.JSONEncodeActionSpec(action, actionSpecEncoder{})
	if err != nil {
		return containerConfig{}, err
	}

	env["NEBULA_ACTION_SPEC"] = encodedSpec

	var joined []string

	for k, v := range env {
		joined = append(joined, fmt.Sprintf("%s=%s", k, v))
	}

	containerCfg := &container.Config{
		Env:   joined,
		Image: action.Image,
	}

	hostCfg := &container.HostConfig{
		Binds: []string{
			e.hostDockerSocketMount(),
			runnerMount,
			buildMount,
		},
		Privileged: true,
		AutoRemove: true,
	}

	networkCfg := &network.NetworkingConfig{}

	return containerConfig{
		container: containerCfg,
		host:      hostCfg,
		network:   networkCfg,
	}, nil
}

func (e *Executor) createContainer(ctx context.Context, r execution.ExecutorRuntime, config containerConfig) (string, errors.Error) {
	response, err := e.client.ContainerCreate(ctx, config.container, config.host, config.network, config.name)
	if err != nil {
		return "", errors.NewDockerContainerCreateError(config.name).WithCause(err)
	}

	r.Logger().Debug("action-container-created", "container", config.name, "container-id", response.ID)

	return response.ID, nil
}

func (e *Executor) setupLogging(ctx context.Context, r execution.ExecutorRuntime, containerID string) errors.Error {
	logOpts := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}

	logs, err := e.client.ContainerLogs(ctx, containerID, logOpts)
	if err != nil {
		return errors.NewDockerContainerLogsError().WithCause(err)
	}

	if _, err := stdcopy.StdCopy(r.IO().Out, r.IO().ErrOut, logs); err != nil {
		return errors.NewDockerContainerLogsError().WithCause(err)
	}

	return nil
}

func (e *Executor) hostDockerSocketMount() string {
	return fmt.Sprintf("%s:%s", e.socketPath, e.socketPath)
}

func NewExecutor(opts ExecutorOptions) (*Executor, errors.Error) {
	c, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, errors.NewDockerClientCreateError().WithCause(err)
	}

	if opts.DockerHostSocketPath == "" {
		opts.DockerHostSocketPath = defaultDockerHostSocketPath
	}

	return &Executor{
		client:     c,
		socketPath: opts.DockerHostSocketPath,
		registryAuth: types.AuthConfig{
			Username:      opts.RegistryUser,
			Password:      opts.RegistryPass,
			ServerAddress: opts.Registry,
		},
	}, nil
}

type actionSpecEncoder struct{}

func (e actionSpecEncoder) Encode(m []byte) (string, errors.Error) {
	encoded := base64.StdEncoding.EncodeToString(m)

	return encoded, nil
}
