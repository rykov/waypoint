// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointDocumentationMapper hashicorp waypoint documentation mapper
//
// swagger:model hashicorp.waypoint.Documentation.Mapper
type HashicorpWaypointDocumentationMapper struct {

	// description
	Description string `json:"description,omitempty"`

	// input
	Input string `json:"input,omitempty"`

	// output
	Output string `json:"output,omitempty"`
}

// Validate validates this hashicorp waypoint documentation mapper
func (m *HashicorpWaypointDocumentationMapper) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp waypoint documentation mapper based on context it is used
func (m *HashicorpWaypointDocumentationMapper) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointDocumentationMapper) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointDocumentationMapper) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointDocumentationMapper
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
