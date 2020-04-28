package evaluate

import "github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/resolve"

type Option func(e *Evaluator)

func WithDataTypeResolver(resolver resolve.DataTypeResolver) Option {
	return func(e *Evaluator) {
		e.dataTypeResolver = resolver
	}
}

func WithSecretTypeResolver(resolver resolve.SecretTypeResolver) Option {
	return func(e *Evaluator) {
		e.secretTypeResolver = resolver
	}
}

func WithConnectionTypeResolver(resolver resolve.ConnectionTypeResolver) Option {
	return func(e *Evaluator) {
		e.connectionTypeResolver = resolver
	}
}

func WithOutputTypeResolver(resolver resolve.OutputTypeResolver) Option {
	return func(e *Evaluator) {
		e.outputTypeResolver = resolver
	}
}

func WithParameterTypeResolver(resolver resolve.ParameterTypeResolver) Option {
	return func(e *Evaluator) {
		e.parameterTypeResolver = resolver
	}
}

func WithAnswerTypeResolver(resolver resolve.AnswerTypeResolver) Option {
	return func(e *Evaluator) {
		e.answerTypeResolver = resolver
	}
}

func WithInvocationResolver(resolver resolve.InvocationResolver) Option {
	return func(e *Evaluator) {
		e.invocationResolver = resolver
	}
}

func WithLanguage(lang Language) Option {
	return func(e *Evaluator) {
		e.lang = lang
	}
}

func WithInvokeFunc(fn InvokeFunc) Option {
	return func(e *Evaluator) {
		e.invoke = fn
	}
}

func WithResultMapper(rm ResultMapper) Option {
	return func(e *Evaluator) {
		e.resultMapper = rm
	}
}
