package fnlib

import (
	"context"
	"reflect"
	"strings"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
)

var concatDescriptor = fn.DescriptorFuncs{
	DescriptionFunc: func() string { return "Concatenates string arguments into a single string" },
	PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
		if len(args) == 0 {
			return fn.StaticInvoker(""), nil
		}

		fn := fn.InvokerFunc(func(ctx context.Context) (m interface{}, err error) {
			strs := make([]string, len(args))
			for i, iarg := range args {
				switch arg := iarg.(type) {
				case string:
					strs[i] = arg
				default:
					return nil, &fn.PositionalArgError{
						Arg: i + 1,
						Cause: &fn.UnexpectedTypeError{
							Wanted: []reflect.Type{reflect.TypeOf("")},
							Got:    reflect.TypeOf(arg),
						},
					}
				}
			}

			return strings.Join(strs, ""), nil
		})
		return fn, nil
	},
}
