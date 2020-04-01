// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// ValueType A value type
//
// swagger:model ValueType
type ValueType string

const (

	// ValueTypeString captures enum value "string"
	ValueTypeString ValueType = "string"
)

// for schema
var valueTypeEnum []interface{}

func init() {
	var res []ValueType
	if err := json.Unmarshal([]byte(`["string"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		valueTypeEnum = append(valueTypeEnum, v)
	}
}

func (m ValueType) validateValueTypeEnum(path, location string, value ValueType) error {
	if err := validate.Enum(path, location, value, valueTypeEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this value type
func (m ValueType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateValueTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
