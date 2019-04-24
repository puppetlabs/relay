// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/swag"
)

// DeleteResponse delete response
// swagger:model DeleteResponse
type DeleteResponse struct {

	// resource id
	ResourceID string `json:"resource_id,omitempty"`

	// success
	Success bool `json:"success,omitempty"`
}

// Validate validates this delete response
func (m *DeleteResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *DeleteResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *DeleteResponse) UnmarshalBinary(b []byte) error {
	var res DeleteResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
