// Package jsonpath implements the template format used by kubectl and adopted
// by Ni.
//
// See https://kubernetes.io/docs/reference/kubectl/jsonpath/ for more
// information.
//
// In some cases, it deviates slightly from the syntax accepted by kubectl's
// JSONPath expressions:
//
// - The use of \ to escape the next character in identifiers is not supported.
// - The use of @['x.y'] (equivalent to @.x.y) inside brackets is not supported,
//   as it could conflict with an actual key in a JSON object.
package jsonpath

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"text/scanner"
	"unicode"

	gval "github.com/puppetlabs/paesslerag-gval"
	jsonpath "github.com/puppetlabs/paesslerag-jsonpath"
)

func ExpressionLanguage() gval.Language {
	return expressionLanguage
}

type TemplateOption func(tl *templateLanguage)

func WithExpressionLanguageVariableVisitor(visitor jsonpath.VariableVisitor) TemplateOption {
	return func(tl *templateLanguage) {
		tl.tpl = gval.NewLanguage(tl.tpl, gval.VariableSelector(jsonpath.VariableSelector(visitor)))
	}
}

func TemplateLanguage(opts ...TemplateOption) gval.Language {
	return newTemplateLanguage(opts).generate()
}

func eval(ctx context.Context, ev gval.Evaluable, parameter interface{}) (string, error) {
	v, err := ev(ctx, parameter)
	if err != nil {
		return "", err
	}

	switch vt := v.(type) {
	case nil:
		return "", nil
	case []interface{}:
		vs := make([]string, len(vt))
		for i, vi := range vt {
			vs[i] = fmt.Sprintf("%v", vi)
		}

		return strings.Join(vs, " "), nil
	default:
		return fmt.Sprintf("%v", vt), nil
	}
}

func concat(ctx context.Context, p *gval.Parser, a gval.Evaluable) (gval.Evaluable, error) {
	b, err := p.ParseExpression(ctx)
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, parameter interface{}) (interface{}, error) {
		ea, err := eval(ctx, a, parameter)
		if err != nil {
			return nil, err
		}

		eb, err := eval(ctx, b, parameter)
		if err != nil {
			return nil, err
		}

		return ea + eb, nil
	}, nil
}

func manyPrefixExtension(runes []rune, ext func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error)) gval.Language {
	l := gval.NewLanguage()
	for _, r := range runes {
		l = gval.NewLanguage(l, gval.PrefixExtension(r, ext))
	}
	return l
}

func parseRange(ctx context.Context, p *gval.Parser, lang gval.Language) (gval.Evaluable, error) {
	query, err := expressionLanguage.ChainEvaluable(ctx, p)
	if err != nil {
		return nil, err
	}

	switch p.Scan() {
	case '}':
	default:
		return nil, p.Expected("JSONPath template range", '}')
	}

	sub, err := lang.ChainEvaluable(ctx, p)
	if err != nil {
		return nil, err
	}

	switch p.Scan() {
	case scanner.Ident:
		if p.TokenText() == "end" {
			break
		}

		fallthrough
	default:
		return nil, p.Expected("JSONPath template range end")
	}

	return func(ctx context.Context, parameter interface{}) (interface{}, error) {
		candidate, err := query(ctx, parameter)
		if err != nil {
			return nil, err
		}

		var s string
		if els, ok := candidate.([]interface{}); ok {
			for _, el := range els {
				v, err := sub.EvalString(ctx, el)
				if err != nil {
					return nil, err
				}

				s += v
			}
		}

		return s, nil
	}, nil
}

func eq(a, b interface{}) bool {
	// Support matrix-y == against scalar values.
	if as, ok := a.([]interface{}); ok {
		for _, av := range as {
			if reflect.DeepEqual(av, b) {
				return true
			}
		}
	} else if bs, ok := b.([]interface{}); ok {
		for _, bv := range bs {
			if reflect.DeepEqual(a, bv) {
				return true
			}
		}
	}

	return reflect.DeepEqual(a, b)
}

// expressionLanguage is the language of JSONPath expressions
var expressionLanguage = gval.NewLanguage(
	gval.Arithmetic(),
	gval.Text(),
	gval.PropositionalLogic(),
	jsonpath.Language(jsonpath.AllowMissingKeys(true)),
	gval.InfixOperator("==", func(a, b interface{}) (interface{}, error) { return eq(a, b), nil }),
	gval.InfixOperator("!=", func(a, b interface{}) (interface{}, error) { return !eq(a, b), nil }),
	gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
		p.SetIsIdentRuneFunc(func(r rune, pos int) bool {
			return unicode.IsLetter(r) ||
				r == '_' ||
				(pos > 0 && (unicode.IsDigit(r) || r == '/' || r == '-'))
		})

		switch p.Scan() {
		case '.', '[':
			// For the first character, we allow omitting the '$' or '@'.
			p.Camouflage("JSONPath expression", '$', '@', '.', '[', '(')
			fallthrough
		case '$', '@':
			// Also, in this specific case, '@' is the same as '$'.
			return jsonpath.Parse(ctx, p, jsonpath.AllowMissingKeys(true))
		default:
			p.Camouflage("JSONPath expression")
			return p.ParseExpression(ctx)
		}
	}),
)

// templateLanguage is the total language, which includes literal handling outside of curly braces
type templateLanguage struct {
	tpl gval.Language
}

func (tl *templateLanguage) generate() gval.Language {
	return gval.Late(func(lang gval.Language) gval.Language {
		tpl := gval.NewLanguage(
			tl.tpl,
			gval.PrefixMetaPrefix(scanner.Ident, func(ctx context.Context, p *gval.Parser) (call string, alternative func() (gval.Evaluable, error), err error) {
				token := p.TokenText()
				return token,
					func() (gval.Evaluable, error) {
						if token == "range" {
							return parseRange(ctx, p, lang)
						}

						p.Camouflage("JSONPath template")
						return p.Const(""), nil
					},
					nil
			}),
		)

		return gval.NewLanguage(
			gval.Init(func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				p.SetWhitespace()
				p.SetMode(scanner.ScanIdents)
				p.SetIsIdentRuneFunc(func(ch rune, i int) bool { return ch > 0 && ch != '{' })

				return p.ParseExpression(ctx)
			}),
			gval.PrefixExtension(scanner.Ident, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				return concat(ctx, p, p.Const(p.TokenText()))
			}),
			gval.PrefixExtension(scanner.EOF, func(ctx context.Context, p *gval.Parser) (gval.Evaluable, error) {
				return p.Const(""), nil
			}),
			gval.PrefixLanguage('{', tpl, func(ctx context.Context, p *gval.Parser, eval gval.Evaluable) (gval.Evaluable, error) {
				switch p.Scan() {
				case '}':
				case scanner.Ident:
					if p.TokenText() == "end" {
						p.Camouflage("JSONPath template", '}')
						return eval, nil
					}

					fallthrough
				default:
					return nil, p.Expected("JSONPath template", '}')
				}

				return concat(ctx, p, eval)
			}),
		)
	})
}

func newTemplateLanguage(opts []TemplateOption) *templateLanguage {
	tl := &templateLanguage{
		tpl: expressionLanguage,
	}
	for _, opt := range opts {
		opt(tl)
	}
	return tl
}
