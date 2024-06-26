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

// HashicorpWaypointConfigSource hashicorp waypoint config source
//
// swagger:model hashicorp.waypoint.ConfigSource
type HashicorpWaypointConfigSource struct {

	// application
	Application *HashicorpWaypointRefApplication `json:"application,omitempty"`

	// config is the configuration for the config source.
	Config map[string]string `json:"config,omitempty"`

	// delete may be set to true on SetConfigSource to delete this config source.
	// This is a field on ConfigSource since there are a variety of ways to
	// identify a ConfigSource. Therefore, the recommend deletion process is
	// to query the ConfigSource using GetConfigSource and then set delete
	// on a return value to ensure the correct value is deleted.
	Delete bool `json:"delete,omitempty"`

	// global
	Global HashicorpWaypointRefGlobal `json:"global,omitempty"`

	// hash is set automatically on write and available on read and is a
	// content hash of the configuration. This can be used to determine
	// uniqueness or changes in the configuration. Setting this value with
	// SetConfigSource has no effect and will be overwritten. Note that this
	// hash may take more into account than just "config" as other fields
	// are introduced to this message type.
	Hash string `json:"hash,omitempty"`

	// project
	Project *HashicorpWaypointRefProject `json:"project,omitempty"`

	// type of the config source. This should match the plugin name.
	Type string `json:"type,omitempty"`

	// workspace, if set, will limit this config source to a specific
	// workspace. This is in addition to the app scoping above. For example,
	// if you specify project scoping above, and set this too, then only
	// matching projects in the matching workspace will have this config var
	// set.
	Workspace *HashicorpWaypointRefWorkspace `json:"workspace,omitempty"`
}

// Validate validates this hashicorp waypoint config source
func (m *HashicorpWaypointConfigSource) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateApplication(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateProject(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateWorkspace(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointConfigSource) validateApplication(formats strfmt.Registry) error {
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

func (m *HashicorpWaypointConfigSource) validateProject(formats strfmt.Registry) error {
	if swag.IsZero(m.Project) { // not required
		return nil
	}

	if m.Project != nil {
		if err := m.Project.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("project")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("project")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointConfigSource) validateWorkspace(formats strfmt.Registry) error {
	if swag.IsZero(m.Workspace) { // not required
		return nil
	}

	if m.Workspace != nil {
		if err := m.Workspace.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("workspace")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("workspace")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this hashicorp waypoint config source based on the context it is used
func (m *HashicorpWaypointConfigSource) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApplication(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateProject(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateWorkspace(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *HashicorpWaypointConfigSource) contextValidateApplication(ctx context.Context, formats strfmt.Registry) error {

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

func (m *HashicorpWaypointConfigSource) contextValidateProject(ctx context.Context, formats strfmt.Registry) error {

	if m.Project != nil {
		if err := m.Project.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("project")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("project")
			}
			return err
		}
	}

	return nil
}

func (m *HashicorpWaypointConfigSource) contextValidateWorkspace(ctx context.Context, formats strfmt.Registry) error {

	if m.Workspace != nil {
		if err := m.Workspace.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("workspace")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("workspace")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointConfigSource) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointConfigSource) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointConfigSource
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
