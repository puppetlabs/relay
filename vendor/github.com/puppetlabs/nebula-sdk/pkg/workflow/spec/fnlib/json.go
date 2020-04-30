package fnlib

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
)

var jsonUnmarshalDescriptor = fn.DescriptorFuncs{
	DescriptionFunc: func() string { return "Unmarshals a JSON-encoded string into the specification" },
	PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
		if len(args) != 1 {
			return nil, &fn.ArityError{Wanted: []int{1}, Got: len(args)}
		}

		fn := fn.InvokerFunc(func(ctx context.Context) (m interface{}, err error) {
			var b []byte

			switch arg := args[0].(type) {
			case string:
				b = []byte(arg)
			default:
				return nil, &fn.PositionalArgError{
					Arg: 1,
					Cause: &fn.UnexpectedTypeError{
						Wanted: []reflect.Type{reflect.TypeOf("")},
						Got:    reflect.TypeOf(arg),
					},
				}
			}

			err = json.Unmarshal(b, &m)
			return
		})
		return fn, nil
	},
}
