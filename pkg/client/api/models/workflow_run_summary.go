// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// WorkflowRunSummary workflow run summary
//
// swagger:model WorkflowRunSummary
type WorkflowRunSummary struct {
	WorkflowRunIdentifier

	// created by
	CreatedBy WorkflowRunCreatedBySummary `json:"created_by,omitempty"`

	// revision
	// Required: true
	Revision *WorkflowRevisionIdentifier `json:"revision"`

	// state
	// Required: true
	State *WorkflowRunStateSummary `json:"state"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *WorkflowRunSummary) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 WorkflowRunIdentifier
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.WorkflowRunIdentifier = aO0

	// AO1
	var dataAO1 struct {
		CreatedBy WorkflowRunCreatedBySummary `json:"created_by,omitempty"`

		Revision *WorkflowRevisionIdentifier `json:"revision"`

		State *WorkflowRunStateSummary `json:"state"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.CreatedBy = dataAO1.CreatedBy

	m.Revision = dataAO1.Revision

	m.State = dataAO1.State

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m WorkflowRunSummary) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.WorkflowRunIdentifier)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		CreatedBy WorkflowRunCreatedBySummary `json:"created_by,omitempty"`

		Revision *WorkflowRevisionIdentifier `json:"revision"`

		State *WorkflowRunStateSummary `json:"state"`
	}

	dataAO1.CreatedBy = m.CreatedBy

	dataAO1.Revision = m.Revision

	dataAO1.State = m.State

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this workflow run summary
func (m *WorkflowRunSummary) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with WorkflowRunIdentifier
	if err := m.WorkflowRunIdentifier.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRevision(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateState(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *WorkflowRunSummary) validateRevision(formats strfmt.Registry) error {

	if err := validate.Required("revision", "body", m.Revision); err != nil {
		return err
	}

	if m.Revision != nil {
		if err := m.Revision.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("revision")
			}
			return err
		}
	}

	return nil
}

func (m *WorkflowRunSummary) validateState(formats strfmt.Registry) error {

	if err := validate.Required("state", "body", m.State); err != nil {
		return err
	}

	if m.State != nil {
		if err := m.State.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("state")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *WorkflowRunSummary) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkflowRunSummary) UnmarshalBinary(b []byte) error {
	var res WorkflowRunSummary
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
