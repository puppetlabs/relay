package resolve

import (
	"context"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fnlib"
	gval "github.com/puppetlabs/paesslerag-gval"
)

type MemoryDataTypeResolver struct {
	m map[string]interface{}
}

var _ DataTypeResolver = &MemoryDataTypeResolver{}

func (mr *MemoryDataTypeResolver) ResolveData(ctx context.Context, query string) (interface{}, error) {
	pl, err := gval.NewLanguage(gval.Base()).NewEvaluable(query)
	if err != nil {
		return "", &DataQueryError{Query: query}
	}

	v, err := pl(ctx, mr.m)
	if err != nil {
		return "", &DataNotFoundError{Query: query}
	}

	return v, nil
}

func NewMemoryDataTypeResolver(m map[string]interface{}) *MemoryDataTypeResolver {
	return &MemoryDataTypeResolver{m: m}
}

type MemorySecretTypeResolver struct {
	m map[string]string
}

var _ SecretTypeResolver = &MemorySecretTypeResolver{}

func (mr *MemorySecretTypeResolver) ResolveSecret(ctx context.Context, name string) (string, error) {
	s, ok := mr.m[name]
	if !ok {
		return "", &SecretNotFoundError{Name: name}
	}

	return s, nil
}

func NewMemorySecretTypeResolver(m map[string]string) *MemorySecretTypeResolver {
	return &MemorySecretTypeResolver{m: m}
}

type MemoryConnectionKey struct {
	Type string
	Name string
}

type MemoryConnectionTypeResolver struct {
	m map[MemoryConnectionKey]interface{}
}

var _ ConnectionTypeResolver = &MemoryConnectionTypeResolver{}

func (mr *MemoryConnectionTypeResolver) ResolveConnection(ctx context.Context, connectionType, name string) (interface{}, error) {
	o, ok := mr.m[MemoryConnectionKey{Type: connectionType, Name: name}]
	if !ok {
		return "", &ConnectionNotFoundError{Type: connectionType, Name: name}
	}

	return o, nil
}

func NewMemoryConnectionTypeResolver(m map[MemoryConnectionKey]interface{}) *MemoryConnectionTypeResolver {
	return &MemoryConnectionTypeResolver{m: m}
}

type MemoryOutputKey struct {
	From string
	Name string
}

type MemoryOutputTypeResolver struct {
	m map[MemoryOutputKey]interface{}
}

var _ OutputTypeResolver = &MemoryOutputTypeResolver{}

func (mr *MemoryOutputTypeResolver) ResolveOutput(ctx context.Context, from, name string) (interface{}, error) {
	o, ok := mr.m[MemoryOutputKey{From: from, Name: name}]
	if !ok {
		return "", &OutputNotFoundError{From: from, Name: name}
	}

	return o, nil
}

func NewMemoryOutputTypeResolver(m map[MemoryOutputKey]interface{}) *MemoryOutputTypeResolver {
	return &MemoryOutputTypeResolver{m: m}
}

type MemoryParameterTypeResolver struct {
	m map[string]interface{}
}

var _ ParameterTypeResolver = &MemoryParameterTypeResolver{}

func (mr *MemoryParameterTypeResolver) ResolveParameter(ctx context.Context, name string) (interface{}, error) {
	p, ok := mr.m[name]
	if !ok {
		return nil, &ParameterNotFoundError{Name: name}
	}

	return p, nil
}

func NewMemoryParameterTypeResolver(m map[string]interface{}) *MemoryParameterTypeResolver {
	return &MemoryParameterTypeResolver{m: m}
}

type MemoryAnswerKey struct {
	AskRef string
	Name   string
}

type MemoryAnswerTypeResolver struct {
	m map[MemoryAnswerKey]interface{}
}

var _ AnswerTypeResolver = &MemoryAnswerTypeResolver{}

func (mr *MemoryAnswerTypeResolver) ResolveAnswer(ctx context.Context, askRef, name string) (interface{}, error) {
	o, ok := mr.m[MemoryAnswerKey{AskRef: askRef, Name: name}]
	if !ok {
		return "", &AnswerNotFoundError{AskRef: askRef, Name: name}
	}

	return o, nil
}

func NewMemoryAnswerTypeResolver(m map[MemoryAnswerKey]interface{}) *MemoryAnswerTypeResolver {
	return &MemoryAnswerTypeResolver{m: m}
}

type MemoryInvocationResolver struct {
	m fn.Map
}

var _ InvocationResolver = &MemoryInvocationResolver{}

func (mr *MemoryInvocationResolver) ResolveInvocationPositional(ctx context.Context, name string, args []interface{}) (fn.Invoker, error) {
	f, err := mr.m.Descriptor(name)
	if err != nil {
		return nil, &FunctionResolutionError{Name: name, Cause: err}
	}

	i, err := f.PositionalInvoker(args)
	if err != nil {
		return nil, &FunctionResolutionError{Name: name, Cause: err}
	}

	// Wrap invoker so we can add the function name to errors produced while
	// invoking.
	wi := fn.InvokerFunc(func(ctx context.Context) (interface{}, error) {
		r, err := i.Invoke(ctx)
		if err != nil {
			return nil, &FunctionResolutionError{Name: name, Cause: err}
		}

		return r, nil
	})
	return wi, nil
}

func (mr *MemoryInvocationResolver) ResolveInvocation(ctx context.Context, name string, args map[string]interface{}) (fn.Invoker, error) {
	f, err := mr.m.Descriptor(name)
	if err != nil {
		return nil, &FunctionResolutionError{Name: name, Cause: err}
	}

	i, err := f.KeywordInvoker(args)
	if err != nil {
		return nil, &FunctionResolutionError{Name: name, Cause: err}
	}

	// Wrap invoker so we can add the function name to errors produced while
	// invoking.
	wi := fn.InvokerFunc(func(ctx context.Context) (interface{}, error) {
		r, err := i.Invoke(ctx)
		if err != nil {
			return nil, &FunctionResolutionError{Name: name, Cause: err}
		}

		return r, nil
	})
	return wi, nil
}

func NewMemoryInvocationResolver(m fn.Map) *MemoryInvocationResolver {
	return &MemoryInvocationResolver{m: m}
}

func NewDefaultMemoryInvocationResolver() *MemoryInvocationResolver {
	return NewMemoryInvocationResolver(fnlib.Library())
}
