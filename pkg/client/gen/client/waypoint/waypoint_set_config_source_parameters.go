// Code generated by go-swagger; DO NOT EDIT.

package waypoint

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"

	"github.com/hashicorp/waypoint/pkg/client/gen/models"
)

// NewWaypointSetConfigSourceParams creates a new WaypointSetConfigSourceParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointSetConfigSourceParams() *WaypointSetConfigSourceParams {
	return &WaypointSetConfigSourceParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointSetConfigSourceParamsWithTimeout creates a new WaypointSetConfigSourceParams object
// with the ability to set a timeout on a request.
func NewWaypointSetConfigSourceParamsWithTimeout(timeout time.Duration) *WaypointSetConfigSourceParams {
	return &WaypointSetConfigSourceParams{
		timeout: timeout,
	}
}

// NewWaypointSetConfigSourceParamsWithContext creates a new WaypointSetConfigSourceParams object
// with the ability to set a context for a request.
func NewWaypointSetConfigSourceParamsWithContext(ctx context.Context) *WaypointSetConfigSourceParams {
	return &WaypointSetConfigSourceParams{
		Context: ctx,
	}
}

// NewWaypointSetConfigSourceParamsWithHTTPClient creates a new WaypointSetConfigSourceParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointSetConfigSourceParamsWithHTTPClient(client *http.Client) *WaypointSetConfigSourceParams {
	return &WaypointSetConfigSourceParams{
		HTTPClient: client,
	}
}

/*
WaypointSetConfigSourceParams contains all the parameters to send to the API endpoint

	for the waypoint set config source operation.

	Typically these are written to a http.Request.
*/
type WaypointSetConfigSourceParams struct {

	// Body.
	Body *models.HashicorpWaypointSetConfigSourceRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint set config source params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointSetConfigSourceParams) WithDefaults() *WaypointSetConfigSourceParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint set config source params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointSetConfigSourceParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) WithTimeout(timeout time.Duration) *WaypointSetConfigSourceParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) WithContext(ctx context.Context) *WaypointSetConfigSourceParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) WithHTTPClient(client *http.Client) *WaypointSetConfigSourceParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) WithBody(body *models.HashicorpWaypointSetConfigSourceRequest) *WaypointSetConfigSourceParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the waypoint set config source params
func (o *WaypointSetConfigSourceParams) SetBody(body *models.HashicorpWaypointSetConfigSourceRequest) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointSetConfigSourceParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
