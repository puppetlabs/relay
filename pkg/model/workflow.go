package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/puppetlabs/relay/pkg/config"
)

type WorkflowSecretEntity struct {
	Secret *WorkflowSecretSummary `json:"secret"`
}

type WorkflowIdentifier struct {
	Name string `json:"name"`
}

type WorkflowSummary struct {
	*WorkflowIdentifier

	Description string `json:"description"`
}

type Workflow struct {
	*WorkflowSummary

	CreatedAt      *time.Time       `json:"created_at"`
	UpdatedAt      *time.Time       `json:"updated_at"`
	LatestRevision *RevisionSummary `json:"latest_revision"`
	State          *WorkflowState   `json:"state,omitempty"`
	MostRecentRun  WorkflowRun      `json:"most_recent_run"`
}

type WorkflowRun struct {
	Revision  RevisionSummary `json:"revision"`
	RunNumber uint            `json:"run_number"`

	// TODO: There's a lot of fields here that are missing. Do we want to backfill them?
}

type WorkflowEntity struct {
	Workflow *Workflow `json:"workflow"`
}

type WorkflowState struct {
	Triggers []*WorkflowTriggerState `json:"triggers"`
}

type WorkflowTriggerState struct {
	Name     string                      `json:"name"`
	Revision *RevisionSummary            `json:"revision"`
	Source   *WorkflowTriggerSourceState `json:"source"`
}

type PushWorkflowTriggerSourceState struct {
	Token string `json:"token"`
}

type ScheduleWorkflowTriggerSourceState struct {
	ScheduledAt string `json:"scheduled_at"`
}

type WebhookWorkflowTriggerSourceState struct {
	Endpoint string `json:"endpoint"`
}

type WorkflowTriggerSourceState struct {
	Type     string                              `json:"type"`
	Status   string                              `json:"status"`
	Push     *PushWorkflowTriggerSourceState     `json:"push,omitempty"`
	Schedule *ScheduleWorkflowTriggerSourceState `json:"schedule,omitempty"`
	Webhook  *WebhookWorkflowTriggerSourceState  `json:"webhook,omitempty"`
}

// What we call 'workflow' to users is really a combination of these two api types.
// This is a departure from the api spec but feels justified?
type WorkflowRevision struct {
	Workflow *Workflow `json:"workflow"`
	Revision *Revision `json:"revision"`
}

func NewWorkflowRevision(workflow *Workflow, revision *Revision) *WorkflowRevision {
	return &WorkflowRevision{
		Workflow: workflow,
		Revision: revision,
	}
}

func (w *WorkflowRevision) Output(cfg *config.Config) {
	if cfg.Out == config.OutputTypeJSON {
		w.OutputJSON()
	}
	// TODO: Text outputter
}

func (w *WorkflowRevision) OutputJSON() {
	jsonBytes, _ := json.MarshalIndent(w, "", "  ")

	fmt.Println(string(jsonBytes))
}
