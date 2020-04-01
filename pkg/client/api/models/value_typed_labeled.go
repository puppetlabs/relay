// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ValueTypedLabeled value typed labeled
//
// swagger:model ValueTypedLabeled
type ValueTypedLabeled struct {
	ValueTyped

	// The label for this value
	Label string `json:"label,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ValueTypedLabeled) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 ValueTyped
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.ValueTyped = aO0

	// AO1
	var dataAO1 struct {
		Label string `json:"label,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Label = dataAO1.Label

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ValueTypedLabeled) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.ValueTyped)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Label string `json:"label,omitempty"`
	}

	dataAO1.Label = m.Label

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this value typed labeled
func (m *ValueTypedLabeled) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with ValueTyped
	if err := m.ValueTyped.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *ValueTypedLabeled) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ValueTypedLabeled) UnmarshalBinary(b []byte) error {
	var res ValueTypedLabeled
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
