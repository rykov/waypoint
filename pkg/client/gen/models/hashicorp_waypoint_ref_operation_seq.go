// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointRefOperationSeq OperationSeq references an operation by sequence number.
//
// swagger:model hashicorp.waypoint.Ref.OperationSeq
type HashicorpWaypointRefOperationSeq struct {

	// application
	Application *HashicorpWaypointRefApplication `json:"application,omitempty"`

	// number
	Number string `json:"number,omitempty"`
}

// Validate validates this hashicorp waypoint ref operation seq
func (m *HashicorpWaypointRefOperationSeq) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateApplication(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointRefOperationSeq) validateApplication(formats strfmt.Registry) error {
	if swag.IsZero(m.Application) { // not required
		return nil
	}

	if m.Application != nil {
		if err := m.Application.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("application")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("application")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint ref operation seq based on the context it is used
func (m *HashicorpWaypointRefOperationSeq) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApplication(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointRefOperationSeq) contextValidateApplication(ctx context.Context, formats strfmt.Registry) error {

	if m.Application != nil {
		if err := m.Application.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("application")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("application")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointRefOperationSeq) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointRefOperationSeq) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointRefOperationSeq
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
