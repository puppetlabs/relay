// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// WorkflowIdentifier A unique identifier for a workflow
//
// swagger:model WorkflowIdentifier
type WorkflowIdentifier struct {

	// name
	// Required: true
	Name WorkflowName `json:"name"`
}

// Validate validates this workflow identifier
func (m *WorkflowIdentifier) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *WorkflowIdentifier) validateName(formats strfmt.Registry) error {

	if err := m.Name.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("name")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *WorkflowIdentifier) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *WorkflowIdentifier) UnmarshalBinary(b []byte) error {
	var res WorkflowIdentifier
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
