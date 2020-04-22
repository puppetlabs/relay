package model

type RevisionSummary struct {
	Id string `json:"id"`
}

type Revision struct {
	Id  string `json:"id"`
	Raw string `json:"raw"`
	// TODO: gonna take a while to fill out the types here
}

type RevisionEntity struct {
	Access   *EntityAccess `json:"access"`
	Revision *Revision     `json:"revision"`
}

// type WorkflowParameters map[string]WorkflowParameter

// type WorkflowParameter struct {
// 	Default     string `json:"default"`
// 	Description string `json:"description,omitempty"`
// }
