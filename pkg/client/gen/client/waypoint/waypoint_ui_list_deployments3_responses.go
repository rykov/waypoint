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

// WaypointUIListDeployments3Reader is a Reader for the WaypointUIListDeployments3 structure.
type WaypointUIListDeployments3Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *WaypointUIListDeployments3Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewWaypointUIListDeployments3OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewWaypointUIListDeployments3Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewWaypointUIListDeployments3OK creates a WaypointUIListDeployments3OK with default headers values
func NewWaypointUIListDeployments3OK() *WaypointUIListDeployments3OK {
	return &WaypointUIListDeployments3OK{}
}

/*
WaypointUIListDeployments3OK describes a response with status code 200, with default header values.

A successful response.
*/
type WaypointUIListDeployments3OK struct {
	Payload *models.HashicorpWaypointUIListDeploymentsResponse
}

// IsSuccess returns true when this waypoint Ui list deployments3 o k response has a 2xx status code
func (o *WaypointUIListDeployments3OK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this waypoint Ui list deployments3 o k response has a 3xx status code
func (o *WaypointUIListDeployments3OK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this waypoint Ui list deployments3 o k response has a 4xx status code
func (o *WaypointUIListDeployments3OK) IsClientError() bool {
	return false
}

// IsServerError returns true when this waypoint Ui list deployments3 o k response has a 5xx status code
func (o *WaypointUIListDeployments3OK) IsServerError() bool {
	return false
}

// IsCode returns true when this waypoint Ui list deployments3 o k response a status code equal to that given
func (o *WaypointUIListDeployments3OK) IsCode(code int) bool {
	return code == 200
}

func (o *WaypointUIListDeployments3OK) Error() string {
	return fmt.Sprintf("[GET /ui/deployments/project/{application.project}][%d] waypointUiListDeployments3OK  %+v", 200, o.Payload)
}

func (o *WaypointUIListDeployments3OK) String() string {
	return fmt.Sprintf("[GET /ui/deployments/project/{application.project}][%d] waypointUiListDeployments3OK  %+v", 200, o.Payload)
}

func (o *WaypointUIListDeployments3OK) GetPayload() *models.HashicorpWaypointUIListDeploymentsResponse {
	return o.Payload
}

func (o *WaypointUIListDeployments3OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.HashicorpWaypointUIListDeploymentsResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewWaypointUIListDeployments3Default creates a WaypointUIListDeployments3Default with default headers values
func NewWaypointUIListDeployments3Default(code int) *WaypointUIListDeployments3Default {
	return &WaypointUIListDeployments3Default{
		_statusCode: code,
	}
}

/*
WaypointUIListDeployments3Default describes a response with status code -1, with default header values.

An unexpected error response.
*/
type WaypointUIListDeployments3Default struct {
	_statusCode int

	Payload *models.GrpcGatewayRuntimeError
}

// Code gets the status code for the waypoint UI list deployments3 default response
func (o *WaypointUIListDeployments3Default) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this waypoint UI list deployments3 default response has a 2xx status code
func (o *WaypointUIListDeployments3Default) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this waypoint UI list deployments3 default response has a 3xx status code
func (o *WaypointUIListDeployments3Default) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this waypoint UI list deployments3 default response has a 4xx status code
func (o *WaypointUIListDeployments3Default) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this waypoint UI list deployments3 default response has a 5xx status code
func (o *WaypointUIListDeployments3Default) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this waypoint UI list deployments3 default response a status code equal to that given
func (o *WaypointUIListDeployments3Default) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *WaypointUIListDeployments3Default) Error() string {
	return fmt.Sprintf("[GET /ui/deployments/project/{application.project}][%d] Waypoint_UI_ListDeployments3 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointUIListDeployments3Default) String() string {
	return fmt.Sprintf("[GET /ui/deployments/project/{application.project}][%d] Waypoint_UI_ListDeployments3 default  %+v", o._statusCode, o.Payload)
}

func (o *WaypointUIListDeployments3Default) GetPayload() *models.GrpcGatewayRuntimeError {
	return o.Payload
}

func (o *WaypointUIListDeployments3Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.GrpcGatewayRuntimeError)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
