// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry hashicorp waypoint get job stream response terminal event table entry
//
// swagger:model hashicorp.waypoint.GetJobStreamResponse.Terminal.Event.TableEntry
type HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry struct {

	// color
	Color string `json:"color,omitempty"`

	// value
	Value string `json:"value,omitempty"`
}

// Validate validates this hashicorp waypoint get job stream response terminal event table entry
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this hashicorp waypoint get job stream response terminal event table entry based on context it is used
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry) UnmarshalBinary(b []byte) error {
	var res HashicorpWaypointGetJobStreamResponseTerminalEventTableEntry
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
