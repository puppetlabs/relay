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

// ContainerWorkflowStep container workflow step
//
// swagger:model ContainerWorkflowStep
type ContainerWorkflowStep struct {
	AnyWorkflowStep

	// Command arguments
	Args []string `json:"args"`

	// Command to issue
	Command string `json:"command,omitempty"`

	// Container image on which step is executed
	// Required: true
	Image *string `json:"image"`

	// Input script to execute
	Input []string `json:"input"`

	// A URL to a script to run
	InputFile string `json:"inputFile,omitempty"`

	// Variable specification data to provide to the container
	Spec interface{} `json:"spec,omitempty"`

	// Type of step
	// Required: true
	// Enum: [container]
	Type *string `json:"type"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ContainerWorkflowStep) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 AnyWorkflowStep
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.AnyWorkflowStep = aO0

	// AO1
	var dataAO1 struct {
		Args []string `json:"args"`

		Command string `json:"command,omitempty"`

		Image *string `json:"image"`

		Input []string `json:"input"`

		InputFile string `json:"inputFile,omitempty"`

		Spec interface{} `json:"spec,omitempty"`

		Type *string `json:"type"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Args = dataAO1.Args

	m.Command = dataAO1.Command

	m.Image = dataAO1.Image

	m.Input = dataAO1.Input

	m.InputFile = dataAO1.InputFile

	m.Spec = dataAO1.Spec

	m.Type = dataAO1.Type

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ContainerWorkflowStep) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.AnyWorkflowStep)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Args []string `json:"args"`

		Command string `json:"command,omitempty"`

		Image *string `json:"image"`

		Input []string `json:"input"`

		InputFile string `json:"inputFile,omitempty"`

		Spec interface{} `json:"spec,omitempty"`

		Type *string `json:"type"`
	}

	dataAO1.Args = m.Args

	dataAO1.Command = m.Command

	dataAO1.Image = m.Image

	dataAO1.Input = m.Input

	dataAO1.InputFile = m.InputFile

	dataAO1.Spec = m.Spec

	dataAO1.Type = m.Type

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this container workflow step
func (m *ContainerWorkflowStep) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with AnyWorkflowStep
	if err := m.AnyWorkflowStep.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateImage(formats); err != nil {
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

func (m *ContainerWorkflowStep) validateImage(formats strfmt.Registry) error {

	if err := validate.Required("image", "body", m.Image); err != nil {
		return err
	}

	return nil
}

var containerWorkflowStepTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["container"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		containerWorkflowStepTypeTypePropEnum = append(containerWorkflowStepTypeTypePropEnum, v)
	}
}

// property enum
func (m *ContainerWorkflowStep) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, containerWorkflowStepTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *ContainerWorkflowStep) validateType(formats strfmt.Registry) error {

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
func (m *ContainerWorkflowStep) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ContainerWorkflowStep) UnmarshalBinary(b []byte) error {
	var res ContainerWorkflowStep
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
