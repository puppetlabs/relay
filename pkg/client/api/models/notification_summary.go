// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NotificationSummary notification summary
//
// swagger:model NotificationSummary
type NotificationSummary struct {
	NotificationIdentifier

	Lifecycle

	// The attributes of this notification
	Attributes []string `json:"attributes"`

	// The type of event that created the event
	// Required: true
	// Enum: [workflow.created workflow.started workflow.failed workflow.cancelled workflow.rejected workflow.succeeded step.approved step.denied invitation.accepted integration.added step.approval integration.reconnection]
	Type *string `json:"type"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *NotificationSummary) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 NotificationIdentifier
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.NotificationIdentifier = aO0

	// AO1
	var aO1 Lifecycle
	if err := swag.ReadJSON(raw, &aO1); err != nil {
		return err
	}
	m.Lifecycle = aO1

	// AO2
	var dataAO2 struct {
		Attributes []string `json:"attributes"`

		Type *string `json:"type"`
	}
	if err := swag.ReadJSON(raw, &dataAO2); err != nil {
		return err
	}

	m.Attributes = dataAO2.Attributes

	m.Type = dataAO2.Type

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m NotificationSummary) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 3)

	aO0, err := swag.WriteJSON(m.NotificationIdentifier)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	aO1, err := swag.WriteJSON(m.Lifecycle)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO1)
	var dataAO2 struct {
		Attributes []string `json:"attributes"`

		Type *string `json:"type"`
	}

	dataAO2.Attributes = m.Attributes

	dataAO2.Type = m.Type

	jsonDataAO2, errAO2 := swag.WriteJSON(dataAO2)
	if errAO2 != nil {
		return nil, errAO2
	}
	_parts = append(_parts, jsonDataAO2)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this notification summary
func (m *NotificationSummary) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with NotificationIdentifier
	if err := m.NotificationIdentifier.Validate(formats); err != nil {
		res = append(res, err)
	}
	// validation for a type composition with Lifecycle
	if err := m.Lifecycle.Validate(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateAttributes(formats); err != nil {
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

var notificationSummaryAttributesItemsEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["actionable","informational"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		notificationSummaryAttributesItemsEnum = append(notificationSummaryAttributesItemsEnum, v)
	}
}

func (m *NotificationSummary) validateAttributesItemsEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, notificationSummaryAttributesItemsEnum); err != nil {
		return err
	}
	return nil
}

func (m *NotificationSummary) validateAttributes(formats strfmt.Registry) error {

	if swag.IsZero(m.Attributes) { // not required
		return nil
	}

	for i := 0; i < len(m.Attributes); i++ {

		// value enum
		if err := m.validateAttributesItemsEnum("attributes"+"."+strconv.Itoa(i), "body", m.Attributes[i]); err != nil {
			return err
		}

	}

	return nil
}

var notificationSummaryTypeTypePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["workflow.created","workflow.started","workflow.failed","workflow.cancelled","workflow.rejected","workflow.succeeded","step.approved","step.denied","invitation.accepted","integration.added","step.approval","integration.reconnection"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		notificationSummaryTypeTypePropEnum = append(notificationSummaryTypeTypePropEnum, v)
	}
}

// property enum
func (m *NotificationSummary) validateTypeEnum(path, location string, value string) error {
	if err := validate.Enum(path, location, value, notificationSummaryTypeTypePropEnum); err != nil {
		return err
	}
	return nil
}

func (m *NotificationSummary) validateType(formats strfmt.Registry) error {

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
func (m *NotificationSummary) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *NotificationSummary) UnmarshalBinary(b []byte) error {
	var res NotificationSummary
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
