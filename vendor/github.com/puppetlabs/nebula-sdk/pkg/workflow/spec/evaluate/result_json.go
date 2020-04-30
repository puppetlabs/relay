package evaluate

import "github.com/puppetlabs/horsehead/v2/encoding/transfer"

type JSONUnresolvableSecretEnvelope struct {
	Name string `json:"name"`
}

type JSONUnresolvableOutputEnvelope struct {
	From string `json:"from"`
	Name string `json:"name"`
}

type JSONUnresolvableParameterEnvelope struct {
	Name string `json:"name"`
}

type JSONUnresolvableAnswerEnvelope struct {
	AskRef string `json:"askRef"`
	Name   string `json:"name"`
}

type JSONUnresolvableInvocationEnvelope struct {
	Name string `json:"name"`
}

type JSONUnresolvableEnvelope struct {
	Secrets     []*JSONUnresolvableSecretEnvelope     `json:"secrets,omitempty"`
	Outputs     []*JSONUnresolvableOutputEnvelope     `json:"outputs,omitempty"`
	Parameters  []*JSONUnresolvableParameterEnvelope  `json:"parameters,omitempty"`
	Answers     []*JSONUnresolvableAnswerEnvelope     `json:"answers,omitempty"`
	Invocations []*JSONUnresolvableInvocationEnvelope `json:"invocations,omitempty"`
}

func NewJSONUnresolvableEnvelope(ur Unresolvable) *JSONUnresolvableEnvelope {
	env := &JSONUnresolvableEnvelope{}

	if len(ur.Secrets) > 0 {
		env.Secrets = make([]*JSONUnresolvableSecretEnvelope, len(ur.Secrets))
		for i, s := range ur.Secrets {
			env.Secrets[i] = &JSONUnresolvableSecretEnvelope{
				Name: s.Name,
			}
		}
	}

	if len(ur.Outputs) > 0 {
		env.Outputs = make([]*JSONUnresolvableOutputEnvelope, len(ur.Outputs))
		for i, o := range ur.Outputs {
			env.Outputs[i] = &JSONUnresolvableOutputEnvelope{
				From: o.From,
				Name: o.Name,
			}
		}
	}

	if len(ur.Parameters) > 0 {
		env.Parameters = make([]*JSONUnresolvableParameterEnvelope, len(ur.Parameters))
		for i, p := range ur.Parameters {
			env.Parameters[i] = &JSONUnresolvableParameterEnvelope{
				Name: p.Name,
			}
		}
	}

	if len(ur.Answers) > 0 {
		env.Answers = make([]*JSONUnresolvableAnswerEnvelope, len(ur.Answers))
		for i, o := range ur.Answers {
			env.Answers[i] = &JSONUnresolvableAnswerEnvelope{
				AskRef: o.AskRef,
				Name:   o.Name,
			}
		}
	}

	if len(ur.Invocations) > 0 {
		env.Invocations = make([]*JSONUnresolvableInvocationEnvelope, len(ur.Invocations))
		for i, call := range ur.Invocations {
			// TODO: Add cause?
			env.Invocations[i] = &JSONUnresolvableInvocationEnvelope{
				Name: call.Name,
			}
		}
	}

	return env
}

type JSONResultEnvelope struct {
	Value        transfer.JSONInterface    `json:"value"`
	Unresolvable *JSONUnresolvableEnvelope `json:"unresolvable"`
	Complete     bool                      `json:"complete"`
}

func NewJSONResultEnvelope(rv *Result) *JSONResultEnvelope {
	return &JSONResultEnvelope{
		Value:        transfer.JSONInterface{Data: rv.Value},
		Unresolvable: NewJSONUnresolvableEnvelope(rv.Unresolvable),
		Complete:     rv.Complete(),
	}
}
