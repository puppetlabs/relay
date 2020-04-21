package model

import "time"

type Workflow struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	// TODO: last run
	// TODO: latest revision
	// TODO: Workflow state
}

type WorkflowEntity struct {
	Access   *EntityAccess `json:"access"`
	Workflow *Workflow     `json:"workflow"`
}
