package client

import (
	"net/http"
	"reflect"

	"github.com/go-openapi/runtime"
	"github.com/puppetlabs/errawr-go/v2/pkg/encoding"
	"github.com/puppetlabs/errawr-go/v2/pkg/errawr"
	"github.com/puppetlabs/nebula-cli/pkg/client/api/models"
	"github.com/puppetlabs/nebula-cli/pkg/errors"
)

func translateModelErrorToEnvelope(merr *models.Error) *encoding.ErrorDisplayEnvelope {
	desc := &encoding.ErrorDescription{}
	if merr.Description != nil {
		desc.Friendly = merr.Description.Friendly
		desc.Technical = merr.Description.Technical
	}

	fmt := &encoding.ErrorDescription{}
	if merr.Formatted != nil {
		fmt.Friendly = merr.Formatted.Friendly
		fmt.Technical = merr.Formatted.Technical
	}

	args, _ := merr.Arguments.(map[string]interface{})

	items := make(map[string]*encoding.ErrorDisplayEnvelope, len(merr.Items))
	for i, item := range merr.Items {
		items[i] = translateModelErrorToEnvelope(&item)
	}

	return &encoding.ErrorDisplayEnvelope{
		Domain:      merr.Domain,
		Section:     merr.Section,
		Code:        merr.Code,
		Title:       merr.Title,
		Description: desc,
		Formatted:   fmt,
		Arguments:   args,
		Items:       items,
	}
}

func translateRuntimeError(err error) errawr.Error {
	if err == nil {
		return nil
	}

	var rv reflect.Value

	type coder interface {
		Code() int
	}

	if c, ok := err.(coder); ok && c.Code() == http.StatusUnauthorized {
		return errors.NewClientNotLoggedIn()
	}

	if rerr, ok := err.(*runtime.APIError); ok {
		rv = reflect.ValueOf(rerr.Response)
	} else {
		rv = reflect.ValueOf(err)
	}

	for rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return errors.NewClientUnexpectedResponseTypeError().WithCause(err)
	}

	pv := rv.FieldByName("Payload")
	if !pv.IsValid() {
		return errors.NewClientUnexpectedResponseTypeError().WithCause(err)
	}
	for pv.Kind() == reflect.Ptr {
		pv = pv.Elem()
	}

	ev := pv.FieldByName("Error")
	if !ev.IsValid() || !ev.CanInterface() {
		return errors.NewClientUnexpectedResponseTypeError().WithCause(err)
	}

	merr, ok := ev.Interface().(*models.Error)
	if !ok || merr == nil {
		return errors.NewClientUnexpectedResponseTypeError().WithCause(err)
	}

	return translateModelErrorToEnvelope(merr).AsError()
}
