package evaluate

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func mapstructureHookFunc(ctx context.Context, e *Evaluator, u *Unresolvable) mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		depth := -1

		// Copy so we can potentially use the zero value below.
		check := to
		for check.Kind() == reflect.Ptr {
			check = check.Elem()
		}

		if check.Kind() == reflect.Struct {
			// We only evaluate one level of nesting for structs, because their
			// children will get correctly traversed once the data exists.
			depth = 1
		}

		r, err := e.evaluate(ctx, data, depth)
		if err != nil {
			return nil, err
		} else if !r.Complete() {
			u.extends(r.Unresolvable)

			// We return the zero value of the type to eliminate confusion.
			return reflect.Zero(to).Interface(), nil
		}

		return r.Value, nil
	}
}
