package model

import "time"

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
	// TODO: last run
}

type WorkflowEntity struct {
	Access   *EntityAccess `json:"access"`
	Workflow *Workflow     `json:"workflow"`
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
