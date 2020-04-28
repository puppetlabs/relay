package fnlib

import (
	"context"
	"reflect"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
)

var appendDescriptor = fn.DescriptorFuncs{
	DescriptionFunc: func() string { return "Adds new items to a given array, returning a new array" },
	PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
		if len(args) < 2 {
			return nil, &fn.ArityError{Wanted: []int{2}, Variadic: true, Got: len(args)}
		}

		fn := fn.InvokerFunc(func(ctx context.Context) (m interface{}, err error) {
			base, ok := args[0].([]interface{})
			if !ok {
				return nil, &fn.PositionalArgError{
					Arg: 1,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{
							reflect.TypeOf([]interface{}(nil)),
						},
						Got: reflect.TypeOf(args[0]),
					},
				}
			}

			new := append([]interface{}{}, base...)
			new = append(new, args[1:]...)
			return new, nil
		})
		return fn, nil
	},
}
