// Code generated by go-swagger; DO NOT EDIT.

package waypoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/waypoint/pkg/client/gen/models"
)

// WaypointRunPipeline2Reader is a Reader for the WaypointRunPipeline2 structure.
type WaypointRunPipeline2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointRunPipeline2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointRunPipeline2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointRunPipeline2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointRunPipeline2OK creates a WaypointRunPipeline2OK with default headers values
func NewWaypointRunPipeline2OK() *WaypointRunPipeline2OK {
	return &WaypointRunPipeline2OK{}
}

/*
WaypointRunPipeline2OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointRunPipeline2OK struct {
	Payload *models.HashicorpWaypointRunPipelineResponse
}

// IsSuccess returns true when this waypoint run pipeline2 o k response has a 2xx status code
func (o *WaypointRunPipeline2OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint run pipeline2 o k response has a 3xx status code
func (o *WaypointRunPipeline2OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint run pipeline2 o k response has a 4xx status code
func (o *WaypointRunPipeline2OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint run pipeline2 o k response has a 5xx status code
func (o *WaypointRunPipeline2OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint run pipeline2 o k response a status code equal to that given
func (o *WaypointRunPipeline2OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointRunPipeline2OK) Error() string {
	return fmt.Sprintf("[POST /pipeline/{pipeline.id}/run][%d] waypointRunPipeline2OK  %+v", 200, o.Payload)
}

func (o *WaypointRunPipeline2OK) String() string {
	return fmt.Sprintf("[POST /pipeline/{pipeline.id}/run][%d] waypointRunPipeline2OK  %+v", 200, o.Payload)
}

func (o *WaypointRunPipeline2OK) GetPayload() *models.HashicorpWaypointRunPipelineResponse {
	return o.Payload
}

func (o *WaypointRunPipeline2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointRunPipelineResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointRunPipeline2Default creates a WaypointRunPipeline2Default with default headers values
func NewWaypointRunPipeline2Default(code int) *WaypointRunPipeline2Default {
	return &WaypointRunPipeline2Default{
		_statusCode: code,
	}
}

/*
WaypointRunPipeline2Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointRunPipeline2Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint run pipeline2 default response
func (o *WaypointRunPipeline2Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint run pipeline2 default response has a 2xx status code
func (o *WaypointRunPipeline2Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint run pipeline2 default response has a 3xx status code
func (o *WaypointRunPipeline2Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint run pipeline2 default response has a 4xx status code
func (o *WaypointRunPipeline2Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint run pipeline2 default response has a 5xx status code
func (o *WaypointRunPipeline2Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint run pipeline2 default response a status code equal to that given
func (o *WaypointRunPipeline2Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointRunPipeline2Default) Error() string {
	return fmt.Sprintf("[POST /pipeline/{pipeline.id}/run][%d] Waypoint_RunPipeline2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointRunPipeline2Default) String() string {
	return fmt.Sprintf("[POST /pipeline/{pipeline.id}/run][%d] Waypoint_RunPipeline2 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointRunPipeline2Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointRunPipeline2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
