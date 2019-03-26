package encoding

import (
	"github.com/puppetlabs/nebula/pkg/errors"
	"github.com/puppetlabs/nebula/pkg/plan/types"
	"sigs.k8s.io/yaml"
)

type ActionSpecEncoder interface {
	Encode(spec []byte) (string, errors.Error)
}

type JSONActionSpecEncoder struct{}

func (e JSONActionSpecEncoder) Encode(spec []byte) (string, errors.Error) {
	b, err := yaml.YAMLToJSON(spec)
	if err != nil {
		return "", errors.NewPlanActionSpecEncodeError().WithCause(err)
	}

	return string(b), nil
}

// JSONActionSpecEncoder encodes action.Spec as JSON, then calls an after encoder if it's
// not nil. If after is nil, then the json string is returned without futher encoding.
func JSONEncodeActionSpec(action *types.Action, after ActionSpecEncoder) (string, errors.Error) {
	e := JSONActionSpecEncoder{}

	s, err := e.Encode(action.Spec)
	if err != nil {
		return "", err
	}

	if after != nil {
		s, err = after.Encode([]byte(s))
		if err != nil {
			return "", err
		}
	}

	return s, nil
}
