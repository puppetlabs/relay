package fn

// Descriptor is a type that describes how a function can be invoked by a caller
type Descriptor interface {
	// Description returns a string that describes what the function does
	Description() string
	// PositionalInvoker takes a slice of values that act like positional arguments
	// to the function. Enforcing order and length constraints is up to the author
	// of the function.
	PositionalInvoker(args []interface{}) (Invoker, error)
	// KeywordInvoker takes its arguments as a map. This acts like labeled or named argments
	// to the function. Enforcing name and length constraints is up to the author
	// of the function.
	KeywordInvoker(args map[string]interface{}) (Invoker, error)
}

// DescriptorFuncs is an adapter that takes anonymous functions that handle methods defined
// in the Descriptor interface. This is a convenience type that allows simple wrapping of
// one-off functions.
type DescriptorFuncs struct {
	DescriptionFunc       func() string
	PositionalInvokerFunc func(args []interface{}) (Invoker, error)
	KeywordInvokerFunc    func(args map[string]interface{}) (Invoker, error)
}

var _ Descriptor = DescriptorFuncs{}

func (df DescriptorFuncs) Description() string {
	if df.DescriptionFunc == nil {
		return "<anonymous>"
	}

	return df.DescriptionFunc()
}

func (df DescriptorFuncs) PositionalInvoker(args []interface{}) (Invoker, error) {
	if df.PositionalInvokerFunc == nil {
		return nil, ErrPositionalArgsNotAccepted
	}

	return df.PositionalInvokerFunc(args)
}

func (df DescriptorFuncs) KeywordInvoker(args map[string]interface{}) (Invoker, error) {
	if df.KeywordInvokerFunc == nil {
		return nil, ErrKeywordArgsNotAccepted
	}

	return df.KeywordInvokerFunc(args)
}

type Map interface {
	Descriptor(name string) (Descriptor, error)
}

type funcMap map[string]Descriptor

func (fm funcMap) Descriptor(name string) (Descriptor, error) {
	fd, found := fm[name]
	if !found {
		return nil, ErrFunctionNotFound
	}

	return fd, nil
}

func NewMap(m map[string]Descriptor) Map {
	return funcMap(m)
}
