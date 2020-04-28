package model

import (
	"time"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/serialize"
)

type RevisionIdentifier struct {
	Id string `json:"id"`
}

type RevisionSummary struct {
	*RevisionIdentifier
}

type Revision struct {
	*RevisionIdentifier

	Parameters WorkflowParameters `json:"parameters"`
	Triggers   []*WorkflowTrigger `json:"triggers"`
	Steps      []*WorkflowStep    `json:"steps"`
	Raw        string             `json:"raw"`
	CreatedAT  *time.Time         `json:"created_at"`
}

type RevisionEntity struct {
	Access   *EntityAccess `json:"access"`
	Revision *Revision     `json:"revision"`
}

type WorkflowParameters map[string]WorkflowParameter

type WorkflowParameter struct {
	Default     string `json:"default"`
	Description string `json:"description,omitempty"`
}

type WorkflowTrigger struct {
	Name    string                  `json:"name"`
	Source  *WorkflowTriggerSource  `json:"source"`
	Binding *WorkflowTriggerBinding `json:"binding"`
}

type WorkflowTriggerSource struct {
	Type                          string `json:"type"`
	PushWorkflowTriggerSource     `json:",inline"`
	ScheduleWorkflowTriggerSource `json:",inline"`
	WebhookWorkflowTriggerSource  `json:",inline"`
}

type PushWorkflowTriggerSource struct {
	Schema map[string]serialize.JSONTree `json:"schema"`
}

type ScheduleWorkflowTriggerSource struct {
	Schedule string `json:"schedule"`
}

type WebhookWorkflowTriggerSource struct {
	ContainerMixin `json:",inline"`
}

type WorkflowTriggerBinding struct {
	Parameters map[string]serialize.JSONTree `json:"parameters,omitempty"`
}

type WorkflowStep struct {
	Name           string                  `json:"name"`
	Description    string                  `json:"description"`
	Type           string                  `json:"type"`
	DependsOn      []string                `json:"depends_on,omitempty"`
	References     *WorkflowDataReferences `json:"references,omitempty"`
	ContainerMixin `json:",inline"`
}

type WorkflowSecretSummary struct {
	Name string `json:"name"`
}

type WorkflowParameterReference struct {
	Name string `json:"name"`
}

type WorkflowOutputReference struct {
	Name string `json:"name"`
	From string `json:"from"`
}

type WorkflowDataReferences struct {
	Secrets    []*WorkflowSecretSummary      `json:"secrets,omitempty"`
	Parameters []*WorkflowParameterReference `json:"parameters,omitempty"`
	Outputs    []*WorkflowOutputReference    `json:"outputs,omitempty"`
}

type ContainerMixin struct {
	Image     string                        `json:"image,omitempty"`
	Spec      map[string]serialize.JSONTree `json:"spec,omitempty"`
	Input     []string                      `json:"input,omitempty"`
	Command   string                        `json:"command,omitempty"`
	Args      []string                      `json:"args,omitempty"`
	InputFile string                        `json:"inputFile,omitempty"`
}
