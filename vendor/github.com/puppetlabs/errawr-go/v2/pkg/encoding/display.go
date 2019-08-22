package encoding

import (
	"fmt"
	"strings"

	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/errawr-go/v2/pkg/impl"
)

type ErrorDisplayEnvelope struct {
	Domain      string                           `json:"domain"`
	Section     string                           `json:"section"`
	Code        string                           `json:"code"`
	Title       string                           `json:"title"`
	Description *ErrorDescription                `json:"description,omitempty"`
	Arguments   map[string]interface{}           `json:"arguments,omitempty"`
	Items       map[string]*ErrorDisplayEnvelope `json:"items,omitempty"`
	Formatted   *ErrorDescription                `json:"formatted,omitempty"`
	Causes      []*ErrorDisplayEnvelope          `json:"causes,omitempty"`
}

func (ede ErrorDisplayEnvelope) AsError() errawr.Error {
	arguments := make(impl.ErrorArguments, len(ede.Arguments))
	for name, argument := range ede.Arguments {
		if argument == nil {
			continue
		}

		arguments[name] = &impl.ErrorArgument{
			Value: argument,
		}
	}

	var items impl.ErrorItems
	if ede.Items != nil {
		items = make(impl.ErrorItems, len(ede.Items))
		for path, item := range ede.Items {
			items[path] = item.AsError()
		}
	}

	prefix := fmt.Sprintf(`%s_%s_`, ede.Domain, ede.Section)
	code := strings.TrimPrefix(ede.Code, prefix)

	description := &impl.ErrorDescription{}
	if ede.Description != nil {
		description.Friendly = ede.Description.Friendly
		description.Technical = ede.Description.Technical
	}

	var e errawr.Error = &impl.Error{
		Version: errawr.Version,
		ErrorDomain: &impl.ErrorDomain{
			Key: ede.Domain,
		},
		ErrorSection: &impl.ErrorSection{
			Key: ede.Section,
		},
		ErrorCode:        code,
		ErrorTitle:       ede.Title,
		ErrorDescription: description,
		ErrorArguments:   arguments,
		ErrorItems:       items,
		ErrorMetadata:    &impl.ErrorMetadata{},
		ErrorSensitivity: errawr.ErrorSensitivityEdge,
	}

	for _, cause := range ede.Causes {
		if cause == nil {
			continue
		}

		e = e.WithCause(cause.AsError())
	}

	return e
}

func ForDisplay(e errawr.Error) *ErrorDisplayEnvelope {
	return ForDisplayWithSensitivity(e, errawr.ErrorSensitivityEdge)
}

func ForDisplayWithSensitivity(e errawr.Error, sensitivity errawr.ErrorSensitivity) *ErrorDisplayEnvelope {
	ede := &ErrorDisplayEnvelope{
		Domain:  e.Domain().Key(),
		Section: e.Section().Key(),
		Code:    e.ID(),
		Title:   e.Title(),
	}

	if items, ok := e.Items(); ok {
		ede.Items = make(map[string]*ErrorDisplayEnvelope, len(items))
		for path, item := range items {
			ede.Items[path] = ForDisplayWithSensitivity(item, sensitivity)
		}
	}

	if e.Sensitivity() > sensitivity {
		return ede
	}

	causes := e.Causes()

	ede.Causes = make([]*ErrorDisplayEnvelope, len(causes))
	for i, cause := range causes {
		ede.Causes[i] = ForDisplayWithSensitivity(cause, sensitivity)
	}

	ede.Description = &ErrorDescription{
		Friendly:  e.Description().Friendly(),
		Technical: e.Description().Technical(),
	}
	ede.Arguments = e.Arguments()
	ede.Formatted = &ErrorDescription{
		Friendly:  e.FormattedDescription().Friendly(),
		Technical: e.FormattedDescription().Technical(),
	}

	return ede
}
