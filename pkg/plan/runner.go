package plan

import (
	"context"

	logging "github.com/puppetlabs/insights-logging"
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/execution/executor"
	"github.com/puppetlabs/nebula/pkg/io"
	"github.com/puppetlabs/nebula/pkg/plan/encoding"
	"github.com/puppetlabs/nebula/pkg/plan/types"
)

type RuntimeFactory interface {
	IO() *io.IO
	Logger() logging.Logger
}

type WorkflowRunner struct {
	wf  *types.Workflow
	arl []*ExecutorWrappedActionRunner
	env map[string]string
}

func (r WorkflowRunner) Name() string {
	return r.wf.Name
}

func (r WorkflowRunner) Run(ctx context.Context, rt RuntimeFactory) errors.Error {
	for _, ar := range r.arl {
		if err := ar.Run(ctx, rt, r.env); err != nil {
			return err
		}
	}

	return nil
}

func NewWorkflowRunnerFromName(name string, p *types.Plan, ex executor.ActionExecutor) (*WorkflowRunner, errors.Error) {
	var (
		workflow *types.Workflow
		runner   = &WorkflowRunner{}
	)

	for _, wf := range p.Workflows {
		if wf.Default && workflow == nil {
			workflow = wf
		}

		if wf.Name == name {
			workflow = wf

			break
		}
	}

	if workflow == nil {
		return nil, errors.NewPlanWorkflowNotFound(name)
	}

	runner.wf = workflow

	actions := make(map[string]*types.Action)
	for _, a := range p.Actions {
		actions[a.Name] = a
	}

	for _, name := range workflow.ActionNames {
		ar := NewExecutorWrappedActionRunner(ex, actions[name])
		runner.arl = append(runner.arl, ar)
	}

	runner.env = make(map[string]string)

	for _, v := range p.Variables {
		runner.env[v.Name] = v.Value
	}

	runner.env["NEBULA_WORKFLOW"] = workflow.Name
	// TODO temporarily faking the run id
	runner.env["NEBULA_RUN_ID"] = "1234"

	return runner, nil
}

type ExecutorWrappedActionRunner struct {
	ex      executor.ActionExecutor
	action  *types.Action
	encoder encoding.ActionSpecEncoder
}

// TODO add env to executor
func (r *ExecutorWrappedActionRunner) Run(ctx context.Context, rt RuntimeFactory, env map[string]string) errors.Error {
	return r.ex.ScheduleAction(ctx, rt, r.action, env)
}

func NewExecutorWrappedActionRunner(ex executor.ActionExecutor, action *types.Action) *ExecutorWrappedActionRunner {
	return &ExecutorWrappedActionRunner{ex: ex, action: action, encoder: encoding.JSONActionSpecEncoder{}}
}
