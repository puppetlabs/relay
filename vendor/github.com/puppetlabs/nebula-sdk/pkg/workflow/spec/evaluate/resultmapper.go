package evaluate

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/puppetlabs/horsehead/v2/encoding/transfer"
)

type ResultMapper interface {
	MapResult(ctx context.Context, r *Result) (*Result, error)
}

type ResultMapperFunc func(ctx context.Context, r *Result) (*Result, error)

var _ ResultMapper = ResultMapperFunc(nil)

func (f ResultMapperFunc) MapResult(ctx context.Context, r *Result) (*Result, error) {
	return f(ctx, r)
}

var IdentityResultMapper = ResultMapperFunc(func(ctx context.Context, r *Result) (*Result, error) {
	return r, nil
})

func ChainResultMappers(rms ...ResultMapper) ResultMapper {
	return ResultMapperFunc(func(ctx context.Context, r *Result) (nr *Result, err error) {
		nr = r
		for _, rm := range rms {
			nr, err = rm.MapResult(ctx, nr)
			if err != nil {
				return
			}
		}
		return
	})
}

type UTF8SafeResultMapper struct{}

func (srm *UTF8SafeResultMapper) do(v interface{}) (interface{}, error) {
	switch vt := v.(type) {
	case []byte:
		return transfer.EncodeJSON(vt)
	case string:
		return transfer.EncodeJSON([]byte(vt))
	case map[string]interface{}:
		r := make(map[string]interface{}, len(vt))
		for k, v := range vt {
			v, err := srm.do(v)
			if err != nil {
				return nil, err
			}

			r[k] = v
		}
		v = r
	case []interface{}:
		s := make([]interface{}, len(vt))
		for i, v := range vt {
			v, err := srm.do(v)
			if err != nil {
				return nil, err
			}

			s[i] = v
		}
		v = s
	}
	return v, nil
}

func (srm *UTF8SafeResultMapper) MapResult(ctx context.Context, r *Result) (*Result, error) {
	v, err := srm.do(r.Value)
	if err != nil {
		return nil, err
	}

	r.Value = v
	return r, nil
}

func NewUTF8SafeResultMapper() *UTF8SafeResultMapper {
	return &UTF8SafeResultMapper{}
}

type JSONResultMapper struct {
	encoder ResultMapper
	indent  int
}

var _ ResultMapper = &JSONResultMapper{}

func (jrm *JSONResultMapper) MapResult(ctx context.Context, r *Result) (*Result, error) {
	r, err := jrm.encoder.MapResult(ctx, r)
	if err != nil {
		return nil, err
	}

	if jrm.indent > 0 {
		r.Value, err = json.MarshalIndent(r.Value, "", strings.Repeat(" ", jrm.indent))
	} else {
		r.Value, err = json.Marshal(r.Value)
	}
	if err != nil {
		return nil, err
	}
	return r, nil
}

type JSONResultMapperOption func(jrm *JSONResultMapper)

func WithJSONResultMapperIndent(spaces int) JSONResultMapperOption {
	return func(jrm *JSONResultMapper) {
		jrm.indent = spaces
	}
}

func NewJSONResultMapper(opts ...JSONResultMapperOption) *JSONResultMapper {
	return &JSONResultMapper{
		encoder: NewUTF8SafeResultMapper(),
	}
}
