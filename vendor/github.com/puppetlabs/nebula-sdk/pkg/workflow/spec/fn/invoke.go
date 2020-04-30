package fn

import "context"

type Invoker interface {
	Invoke(ctx context.Context) (interface{}, error)
}

type InvokerFunc func(ctx context.Context) (interface{}, error)

var _ Invoker = InvokerFunc(nil)

func (fn InvokerFunc) Invoke(ctx context.Context) (interface{}, error) {
	return fn(ctx)
}

func StaticInvoker(value interface{}) Invoker {
	return InvokerFunc(func(_ context.Context) (interface{}, error) { return value, nil })
}
