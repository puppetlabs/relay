package evaluate

import (
	"context"

	gval "github.com/puppetlabs/paesslerag-gval"
	jsonpath "github.com/puppetlabs/paesslerag-jsonpath"
)

func variableSelector(e *Evaluator, r *Result) func(path gval.Evaluables) gval.Evaluable {
	visitor := variableVisitor(e, r)

	return func(path gval.Evaluables) gval.Evaluable {
		return func(ctx context.Context, v interface{}) (interface{}, error) {
			for _, key := range path {
				key, err := key(ctx, v)
				if err != nil {
					return nil, err
				}

				// For consistency we use the JSONPath visitor here even though
				// it isn't strictly necessary.
				var nv interface{}
				err = visitor.VisitChild(ctx, v, key, func(ctx context.Context, pv jsonpath.PathValue) error {
					nv = pv.Value
					return nil
				})
				if err != nil {
					return nil, err
				} else if nv == nil {
					return nil, nil
				}

				v = nv
			}

			return v, nil
		}
	}
}

func variableVisitor(e *Evaluator, r *Result) jsonpath.VariableVisitor {
	return jsonpath.VariableVisitorFuncs{
		VisitChildFunc: func(ctx context.Context, parameter interface{}, key interface{}, next func(ctx context.Context, pv jsonpath.PathValue) error) error {
			return jsonpath.DefaultVariableVisitor().VisitChild(ctx, parameter, key, func(ctx context.Context, pv jsonpath.PathValue) error {
				// Expand just this value without recursing.
				nr, err := e.evaluate(ctx, pv.Value, 1)
				if err != nil {
					return err
				} else if !nr.Complete() {
					r.extends(nr)
					return nil
				}

				pv.Value = nr.Value
				return next(ctx, pv)
			})
		},
	}
}
