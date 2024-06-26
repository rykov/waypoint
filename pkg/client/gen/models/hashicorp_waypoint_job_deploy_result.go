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

// HashicorpWaypointJobDeployResult hashicorp waypoint job deploy result
//
// swagger:model hashicorp.waypoint.Job.DeployResult
type HashicorpWaypointJobDeployResult struct {

	// deployment
	Deployment *HashicorpWaypointDeployment `json:"deployment,omitempty"`
}

// Validate validates this hashicorp waypoint job deploy result
func (m *HashicorpWaypointJobDeployResult) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDeployment(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobDeployResult) validateDeployment(formats strfmt.Registry) error {
	if swag.IsZero(m.Deployment) { // not required
		return nil
	}

	if m.Deployment != nil {
		if err := m.Deployment.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("deployment")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("deployment")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint job deploy result based on the context it is used
func (m *HashicorpWaypointJobDeployResult) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDeployment(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointJobDeployResult) contextValidateDeployment(ctx context.Context, formats strfmt.Registry) error {

	if m.Deployment != nil {
		if err := m.Deployment.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("deployment")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("deployment")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointJobDeployResult) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointJobDeployResult) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointJobDeployResult
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
