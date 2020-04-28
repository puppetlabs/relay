package fnlib

import (
	"context"
	"fmt"
	"reflect"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/fn"
)

func merge(dst, src map[string]interface{}, deep bool) {
	for k, v := range src {
		if deep {
			if dm, ok := dst[k].(map[string]interface{}); ok {
				if sm, ok := v.(map[string]interface{}); ok {
					merge(dm, sm, deep)
					continue
				}
			}
		}

		dst[k] = v
	}
}

func mergeCast(os []interface{}) ([]map[string]interface{}, error) {
	objs := make([]map[string]interface{}, len(os))
	for i, o := range os {
		obj, ok := o.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("array index %d: %+v", i, &fn.UnexpectedTypeError{
				Wanted: []reflect.Type{reflect.TypeOf(map[string]interface{}(nil))},
				Got:    reflect.TypeOf(o),
			})
		}

		objs[i] = obj
	}
	return objs, nil
}

var mergeDescriptor = fn.DescriptorFuncs{
	DescriptionFunc: func() string {
		return `Merges a series of objects, with each object overwriting prior entries.

Merges are performed deeply by default. Use the keyword form and set mode: shallow to perform a shallow merge.`
	},
	PositionalInvokerFunc: func(args []interface{}) (fn.Invoker, error) {
		if len(args) == 0 {
			return fn.StaticInvoker(map[string]interface{}{}), nil
		}

		return fn.InvokerFunc(func(ctx context.Context) (interface{}, error) {
			objs, err := mergeCast(args)
			if err != nil {
				return nil, &fn.PositionalArgError{
					Arg:   1,
					Cause: err,
				}
			}

			r := make(map[string]interface{})
			for _, obj := range objs {
				merge(r, obj, true)
			}
			return r, nil
		}), nil
	},
	KeywordInvokerFunc: func(args map[string]interface{}) (fn.Invoker, error) {
		oi, found := args["objects"]
		if !found {
			return nil, &fn.KeywordArgError{Arg: "objects", Cause: fn.ErrArgNotFound}
		}

		mode, found := args["mode"]
		if !found {
			mode = "deep"
		}

		return fn.InvokerFunc(func(ctx context.Context) (interface{}, error) {
			var deep bool
			switch mode {
			case "deep":
				deep = true
			case "shallow":
				deep = false
			default:
				return nil, &fn.KeywordArgError{
					Arg:   "mode",
					Cause: fmt.Errorf(`unexpected value %q, wanted one of "deep" or "shallow"`, mode),
				}
			}

			os, ok := oi.([]interface{})
			if !ok {

			}

			objs, err := mergeCast(os)
			if err != nil {
				return nil, &fn.KeywordArgError{
					Arg:   "objects",
					Cause: err,
				}
			}

			r := make(map[string]interface{})
			for _, obj := range objs {
				merge(r, obj, deep)
			}
			return r, nil
		}), nil
	},
}
