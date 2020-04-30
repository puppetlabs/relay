package evaluate

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
	"github.com/puppetlabs/nebula-sdk/pkg/util/jsonpath"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/parse"
	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/resolve"
	gval "github.com/puppetlabs/paesslerag-gval"
	gvaljsonpath "github.com/puppetlabs/paesslerag-jsonpath"
)

type Language int

const (
	LanguagePath Language = 1 + iota
	LanguageJSONPath
	LanguageJSONPathTemplate
)

type InvokeFunc func(ctx context.Context, i fn.Invoker) (interface{}, error)

type Evaluator struct {
	lang                   Language
	invoke                 InvokeFunc
	resultMapper           ResultMapper
	dataTypeResolver       resolve.DataTypeResolver
	secretTypeResolver     resolve.SecretTypeResolver
	connectionTypeResolver resolve.ConnectionTypeResolver
	outputTypeResolver     resolve.OutputTypeResolver
	parameterTypeResolver  resolve.ParameterTypeResolver
	answerTypeResolver     resolve.AnswerTypeResolver
	invocationResolver     resolve.InvocationResolver
}

func (e *Evaluator) ScopeTo(tree parse.Tree) *ScopedEvaluator {
	return &ScopedEvaluator{parent: e, tree: tree}
}

func (e *Evaluator) Copy(opts ...Option) *Evaluator {
	if len(opts) == 0 {
		return e
	}

	ne := &Evaluator{}
	*ne = *e

	for _, opt := range opts {
		opt(ne)
	}

	return ne
}

func (e *Evaluator) Evaluate(ctx context.Context, tree parse.Tree, depth int) (*Result, error) {
	r, err := e.evaluate(ctx, tree, depth)
	if err != nil {
		return nil, err
	}

	return e.resultMapper.MapResult(ctx, r)
}

func (e *Evaluator) EvaluateAll(ctx context.Context, tree parse.Tree) (*Result, error) {
	return e.Evaluate(ctx, tree, -1)
}

func (e *Evaluator) EvaluateInto(ctx context.Context, tree parse.Tree, target interface{}) (Unresolvable, error) {
	var u Unresolvable

	d, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructureHookFunc(ctx, e, &u),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToTimeHookFunc(time.RFC3339Nano),
		),
		ZeroFields: true,
		Result:     target,
		TagName:    "spec",
	})
	if err != nil {
		return u, err
	}

	return u, d.Decode(tree)
}

func (e *Evaluator) EvaluateQuery(ctx context.Context, tree parse.Tree, query string) (*Result, error) {
	r := &Result{}

	var pl gval.Language
	switch e.lang {
	case LanguagePath:
		pl = gval.NewLanguage(
			gval.Base(),
			gval.VariableSelector(variableSelector(e, r)),
		)
	case LanguageJSONPath:
		pl = gval.NewLanguage(
			jsonpath.ExpressionLanguage(),
			gval.VariableSelector(gvaljsonpath.VariableSelector(variableVisitor(e, r))),
		)
	case LanguageJSONPathTemplate:
		pl = jsonpath.TemplateLanguage(jsonpath.WithExpressionLanguageVariableVisitor(variableVisitor(e, r)))
	default:
		return nil, ErrUnsupportedLanguage
	}

	path, err := pl.NewEvaluable(query)
	if err != nil {
		return nil, err
	}

	v, err := path(ctx, tree)
	if err != nil {
		return nil, err
	}

	er, err := e.evaluate(ctx, v, -1)
	if err != nil {
		return nil, err
	}

	// Add any other unresolved paths in here (provided by the variable selector).
	er.extends(r)

	return e.resultMapper.MapResult(ctx, er)
}

func (e *Evaluator) evaluateType(ctx context.Context, tm map[string]interface{}) (*Result, error) {
	switch tm["$type"] {
	case "Data":
		query, ok := tm["query"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Data", Cause: &FieldNotFoundError{Name: "query"}}
		}

		value, err := e.dataTypeResolver.ResolveData(ctx, query)
		if serr, ok := err.(*resolve.DataNotFoundError); ok {
			return &Result{
				Value: tm,
				Unresolvable: Unresolvable{Data: []UnresolvableData{
					{Query: serr.Query},
				}},
			}, nil
		} else if err != nil {
			return nil, &InvalidTypeError{Type: "Data", Cause: err}
		}

		return &Result{Value: value}, nil
	case "Secret":
		name, ok := tm["name"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Secret", Cause: &FieldNotFoundError{Name: "name"}}
		}

		value, err := e.secretTypeResolver.ResolveSecret(ctx, name)
		if serr, ok := err.(*resolve.SecretNotFoundError); ok {
			return &Result{
				Value: tm,
				Unresolvable: Unresolvable{Secrets: []UnresolvableSecret{
					{Name: serr.Name},
				}},
			}, nil
		} else if err != nil {
			return nil, &InvalidTypeError{Type: "Secret", Cause: err}
		}

		return &Result{Value: value}, nil
	case "Connection":
		connectionType, ok := tm["type"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Connection", Cause: &FieldNotFoundError{Name: "type"}}
		}

		name, ok := tm["name"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Connection", Cause: &FieldNotFoundError{Name: "name"}}
		}

		value, err := e.connectionTypeResolver.ResolveConnection(ctx, connectionType, name)
		if oerr, ok := err.(*resolve.ConnectionNotFoundError); ok {
			return &Result{
				Value: tm,
				Unresolvable: Unresolvable{Connections: []UnresolvableConnection{
					{Type: oerr.Type, Name: oerr.Name},
				}},
			}, nil
		} else if err != nil {
			return nil, &InvalidTypeError{Type: "Connection", Cause: err}
		}

		return &Result{Value: value}, nil
	case "Output":
		from, ok := tm["from"].(string)
		if !ok {
			// Fall back to old syntax.
			//
			// TODO: Remove this in a second version.
			from, ok = tm["taskName"].(string)
			if !ok {
				return nil, &InvalidTypeError{Type: "Output", Cause: &FieldNotFoundError{Name: "from"}}
			}
		}

		name, ok := tm["name"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Output", Cause: &FieldNotFoundError{Name: "name"}}
		}

		value, err := e.outputTypeResolver.ResolveOutput(ctx, from, name)
		if oerr, ok := err.(*resolve.OutputNotFoundError); ok {
			return &Result{
				Value: tm,
				Unresolvable: Unresolvable{Outputs: []UnresolvableOutput{
					{From: oerr.From, Name: oerr.Name},
				}},
			}, nil
		} else if err != nil {
			return nil, &InvalidTypeError{Type: "Output", Cause: err}
		}

		return &Result{Value: value}, nil
	case "Parameter":
		name, ok := tm["name"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Parameter", Cause: &FieldNotFoundError{Name: "name"}}
		}

		value, err := e.parameterTypeResolver.ResolveParameter(ctx, name)
		if perr, ok := err.(*resolve.ParameterNotFoundError); ok {
			return &Result{
				Value: tm,
				Unresolvable: Unresolvable{Parameters: []UnresolvableParameter{
					{Name: perr.Name},
				}},
			}, nil
		} else if err != nil {
			return nil, &InvalidTypeError{Type: "Parameter", Cause: err}
		}

		return &Result{Value: value}, nil
	case "Answer":
		askRef, ok := tm["askRef"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Answer", Cause: &FieldNotFoundError{Name: "askRef"}}
		}

		name, ok := tm["name"].(string)
		if !ok {
			return nil, &InvalidTypeError{Type: "Answer", Cause: &FieldNotFoundError{Name: "name"}}
		}

		value, err := e.answerTypeResolver.ResolveAnswer(ctx, askRef, name)
		if oerr, ok := err.(*resolve.AnswerNotFoundError); ok {
			return &Result{
				Value: tm,
				Unresolvable: Unresolvable{Answers: []UnresolvableAnswer{
					{AskRef: oerr.AskRef, Name: oerr.Name},
				}},
			}, nil
		} else if err != nil {
			return nil, &InvalidTypeError{Type: "Answer", Cause: err}
		}

		return &Result{Value: value}, nil
	default:
		return &Result{Value: tm}, nil
	}
}

func (e *Evaluator) evaluateEncoding(ctx context.Context, em map[string]interface{}) (*Result, error) {
	ty, ok := em["$encoding"].(string)
	if !ok {
		return &Result{Value: em}, nil
	}

	dr, err := e.evaluate(ctx, em["data"], -1)
	if err != nil {
		return nil, &InvalidEncodingError{Type: ty, Cause: err}
	} else if !dr.Complete() {
		r := &Result{
			Value: map[string]interface{}{
				"$encoding": ty,
				"data":      dr.Value,
			},
		}
		r.extends(dr)
		return r, nil
	}

	data, ok := dr.Value.(string)
	if !ok {
		return nil, &InvalidEncodingError{
			Type: ty,
			Cause: &fn.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf("")},
				Got:    reflect.TypeOf(dr.Value),
			},
		}
	}

	decoded, err := transfer.JSON{
		EncodingType: transfer.EncodingType(ty),
		Data:         data,
	}.Decode()
	if err != nil {
		return nil, &InvalidEncodingError{Type: ty, Cause: err}
	}

	return &Result{Value: string(decoded)}, nil
}

func (e *Evaluator) evaluateInvocation(ctx context.Context, im map[string]interface{}) (*Result, error) {
	var key string
	var value interface{}
	for key, value = range im {
	}

	name := strings.TrimPrefix(key, "$fn.")

	a, err := e.evaluate(ctx, value, -1)
	if err != nil {
		return nil, err
	} else if !a.Complete() {
		r := &Result{
			Value: map[string]interface{}{
				key: a.Value,
			},
		}
		r.extends(a)
		return r, nil
	}

	var invoker fn.Invoker
	switch args := a.Value.(type) {
	case []interface{}:
		invoker, err = e.invocationResolver.ResolveInvocationPositional(ctx, name, args)
	case map[string]interface{}:
		invoker, err = e.invocationResolver.ResolveInvocation(ctx, name, args)
	default:
		invoker, err = e.invocationResolver.ResolveInvocationPositional(ctx, name, []interface{}{args})
	}
	if ierr, ok := err.(*resolve.FunctionResolutionError); ok {
		return &Result{
			Value: im,
			Unresolvable: Unresolvable{Invocations: []UnresolvableInvocation{
				{Name: ierr.Name, Cause: ierr.Cause},
			}},
		}, nil
	} else if err != nil {
		return nil, err
	}

	v, err := e.invoke(ctx, invoker)
	if err != nil {
		return nil, err
	}

	return &Result{Value: v}, nil
}

func (e *Evaluator) evaluate(ctx context.Context, v interface{}, depth int) (*Result, error) {
	if depth == 0 {
		return &Result{Value: v}, nil
	}

	switch vt := v.(type) {
	case []interface{}:
		if depth == 1 {
			return &Result{Value: v}, nil
		}

		r := &Result{}
		l := make([]interface{}, len(vt))
		for i, v := range vt {
			nv, err := e.evaluate(ctx, v, depth-1)
			if err != nil {
				return nil, &PathEvaluationError{
					Path:  strconv.Itoa(i),
					Cause: err,
				}
			}

			r.extends(nv)
			l[i] = nv.Value
		}

		r.Value = l
		return r, nil
	case map[string]interface{}:
		if _, ok := vt["$type"]; ok {
			return e.evaluateType(ctx, vt)
		} else if _, ok := vt["$encoding"]; ok {
			return e.evaluateEncoding(ctx, vt)
		} else if len(vt) == 1 {
			var first string
			for first = range vt {
			}

			if strings.HasPrefix(first, "$fn.") {
				return e.evaluateInvocation(ctx, vt)
			}
		} else if depth == 1 {
			return &Result{Value: v}, nil
		}

		r := &Result{}
		m := make(map[string]interface{}, len(vt))
		for k, v := range vt {
			nv, err := e.evaluate(ctx, v, depth-1)
			if err != nil {
				return nil, &PathEvaluationError{Path: k, Cause: err}
			}

			r.extends(nv)
			m[k] = nv.Value
		}

		r.Value = m
		return r, nil
	default:
		return &Result{Value: v}, nil
	}
}

func NewEvaluator(opts ...Option) *Evaluator {
	e := &Evaluator{
		lang:                   LanguagePath,
		invoke:                 func(ctx context.Context, i fn.Invoker) (interface{}, error) { return i.Invoke(ctx) },
		resultMapper:           IdentityResultMapper,
		dataTypeResolver:       resolve.NoOpDataTypeResolver,
		secretTypeResolver:     resolve.NoOpSecretTypeResolver,
		connectionTypeResolver: resolve.NoOpConnectionTypeResolver,
		outputTypeResolver:     resolve.NoOpOutputTypeResolver,
		parameterTypeResolver:  resolve.NoOpParameterTypeResolver,
		answerTypeResolver:     resolve.NoOpAnswerTypeResolver,
		invocationResolver:     resolve.NewDefaultMemoryInvocationResolver(),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}
