// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ProviderAuthOAuth Metadata for authorizing an OAuth-authenticated integration provider
//
// swagger:model ProviderAuthOAuth
type ProviderAuthOAuth struct {

	// The time at which the redirect URL will no longer be valid
	// Required: true
	// Format: date-time
	ExpiresAt *strfmt.DateTime `json:"expires_at"`

	// A URL at which the integration can be authenticated
	// Required: true
	RedirectURL *string `json:"redirect_url"`

	// type
	// Required: true
	// Enum: [oauth]
	Type *string `json:"type"`
}

// Validate validates this provider auth o auth
func (m *ProviderAuthOAuth) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateExpiresAt(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRedirectURL(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProviderAuthOAuth) validateExpiresAt(formats strfmt.Registry) error {

	if err := validate.Required("expires_at", "body", m.ExpiresAt); err != nil {
		return err
	}

	if err := validate.FormatOf("expires_at", "body", "date-time", m.ExpiresAt.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *ProviderAuthOAuth) validateRedirectURL(formats strfmt.Registry) error {

	if err := validate.Required("redirect_url", "body", m.RedirectURL); err != nil {
		return err
	}

	return nil
}

var providerAuthOAuthTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["oauth"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		providerAuthOAuthTypeTypePropEnum = append(providerAuthOAuthTypeTypePropEnum, v)
	}
}

const (

	// ProviderAuthOAuthTypeOauth captures enum value "oauth"
	ProviderAuthOAuthTypeOauth string = "oauth"
)

// prop value enum
func (m *ProviderAuthOAuth) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, providerAuthOAuthTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ProviderAuthOAuth) validateType(formats strfmt.Registry) error {

	if err := validate.Required("type", "body", m.Type); err != nil {
		return err
	}

	// value enum
	if err := m.validateTypeEnum("type", "body", *m.Type); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProviderAuthOAuth) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProviderAuthOAuth) UnmarshalBinary(b []byte) error {
	var res ProviderAuthOAuth
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
