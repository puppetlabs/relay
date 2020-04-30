package evaluate

import (
	"context"

	"github.com/puppetlabs/nebula-sdk/pkg/workflow/spec/parse"
)

type ScopedEvaluator struct {
	parent *Evaluator
	tree   parse.Tree
}

func (se *ScopedEvaluator) Evaluate(ctx context.Context, depth int) (*Result, error) {
	return se.parent.Evaluate(ctx, se.tree, depth)
}

func (se *ScopedEvaluator) EvaluateAll(ctx context.Context) (*Result, error) {
	return se.parent.EvaluateAll(ctx, se.tree)
}

func (se *ScopedEvaluator) EvaluateInto(ctx context.Context, target interface{}) (Unresolvable, error) {
	return se.parent.EvaluateInto(ctx, se.tree, target)
}

func (se *ScopedEvaluator) EvaluateQuery(ctx context.Context, query string) (*Result, error) {
	return se.parent.EvaluateQuery(ctx, se.tree, query)
}

func (se *ScopedEvaluator) Copy(opts ...Option) *ScopedEvaluator {
	return &ScopedEvaluator{parent: se.parent.Copy(opts...), tree: se.tree}
}
