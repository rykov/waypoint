// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// HashicorpWaypointStatus Status represents the status of an async operation.
//
// swagger:model hashicorp.waypoint.Status
type HashicorpWaypointStatus struct {

	// complete_time is the time the operation completed (success or fail).
	// Format: date-time
	CompleteTime strfmt.DateTime `json:"complete_time,omitempty"`

	// details may be non-empty to provide human-friendly information
	// about the current status. This may change between status updates
	// for the same state to provide updated details about the state.
	Details string `json:"details,omitempty"`

	// error is set if the state == ERROR with the error that occurred.
	Error *GoogleRPCStatus `json:"error,omitempty"`

	// start_time is the time the operation was started.
	// Format: date-time
	StartTime strfmt.DateTime `json:"start_time,omitempty"`

	// state is the state of this operation.
	State *HashicorpWaypointStatusState `json:"state,omitempty"`
}

// Validate validates this hashicorp waypoint status
func (m *HashicorpWaypointStatus) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCompleteTime(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateError(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStartTime(formats); err != nil {
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

func (m *HashicorpWaypointStatus) validateCompleteTime(formats strfmt.Registry) error {
	if swag.IsZero(m.CompleteTime) { // not required
		return nil
	}

	if err := validate.FormatOf("complete_time", "body", "date-time", m.CompleteTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *HashicorpWaypointStatus) validateError(formats strfmt.Registry) error {
	if swag.IsZero(m.Error) { // not required
		return nil
	}

	if m.Error != nil {
		if err := m.Error.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("error")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("error")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointStatus) validateStartTime(formats strfmt.Registry) error {
	if swag.IsZero(m.StartTime) { // not required
		return nil
	}

	if err := validate.FormatOf("start_time", "body", "date-time", m.StartTime.String(), formats); err != nil {
		return err
	}

	return nil
}

func (m *HashicorpWaypointStatus) validateState(formats strfmt.Registry) error {
	if swag.IsZero(m.State) { // not required
		return nil
	}

	if m.State != nil {
		if err := m.State.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("state")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint status based on the context it is used
func (m *HashicorpWaypointStatus) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateError(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateState(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointStatus) contextValidateError(ctx context.Context, formats strfmt.Registry) error {

	if m.Error != nil {
		if err := m.Error.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("error")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("error")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointStatus) contextValidateState(ctx context.Context, formats strfmt.Registry) error {

	if m.State != nil {
		if err := m.State.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("state")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("state")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointStatus) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointStatus) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointStatus
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
