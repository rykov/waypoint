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
)

// NewWaypointListWorkspacesParams creates a new WaypointListWorkspacesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewWaypointListWorkspacesParams() *WaypointListWorkspacesParams {
	return &WaypointListWorkspacesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewWaypointListWorkspacesParamsWithTimeout creates a new WaypointListWorkspacesParams object
// with the ability to set a timeout on a request.
func NewWaypointListWorkspacesParamsWithTimeout(timeout time.Duration) *WaypointListWorkspacesParams {
	return &WaypointListWorkspacesParams{
		timeout: timeout,
	}
}

// NewWaypointListWorkspacesParamsWithContext creates a new WaypointListWorkspacesParams object
// with the ability to set a context for a request.
func NewWaypointListWorkspacesParamsWithContext(ctx context.Context) *WaypointListWorkspacesParams {
	return &WaypointListWorkspacesParams{
		Context: ctx,
	}
}

// NewWaypointListWorkspacesParamsWithHTTPClient creates a new WaypointListWorkspacesParams object
// with the ability to set a custom HTTPClient for a request.
func NewWaypointListWorkspacesParamsWithHTTPClient(client *http.Client) *WaypointListWorkspacesParams {
	return &WaypointListWorkspacesParams{
		HTTPClient: client,
	}
}

/*
WaypointListWorkspacesParams contains all the parameters to send to the API endpoint

	for the waypoint list workspaces operation.

	Typically these are written to a http.Request.
*/
type WaypointListWorkspacesParams struct {

	// ApplicationApplication.
	ApplicationApplication *string

	// ApplicationProject.
	ApplicationProject *string

	// ProjectProject.
	ProjectProject *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the waypoint list workspaces params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointListWorkspacesParams) WithDefaults() *WaypointListWorkspacesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the waypoint list workspaces params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *WaypointListWorkspacesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) WithTimeout(timeout time.Duration) *WaypointListWorkspacesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) WithContext(ctx context.Context) *WaypointListWorkspacesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) WithHTTPClient(client *http.Client) *WaypointListWorkspacesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithApplicationApplication adds the applicationApplication to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) WithApplicationApplication(applicationApplication *string) *WaypointListWorkspacesParams {
	o.SetApplicationApplication(applicationApplication)
	return o
}

// SetApplicationApplication adds the applicationApplication to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) SetApplicationApplication(applicationApplication *string) {
	o.ApplicationApplication = applicationApplication
}

// WithApplicationProject adds the applicationProject to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) WithApplicationProject(applicationProject *string) *WaypointListWorkspacesParams {
	o.SetApplicationProject(applicationProject)
	return o
}

// SetApplicationProject adds the applicationProject to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) SetApplicationProject(applicationProject *string) {
	o.ApplicationProject = applicationProject
}

// WithProjectProject adds the projectProject to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) WithProjectProject(projectProject *string) *WaypointListWorkspacesParams {
	o.SetProjectProject(projectProject)
	return o
}

// SetProjectProject adds the projectProject to the waypoint list workspaces params
func (o *WaypointListWorkspacesParams) SetProjectProject(projectProject *string) {
	o.ProjectProject = projectProject
}

// WriteToRequest writes these params to a swagger request
func (o *WaypointListWorkspacesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.ApplicationApplication != nil {

		// query param application.application
		var qrApplicationApplication string

		if o.ApplicationApplication != nil {
			qrApplicationApplication = *o.ApplicationApplication
		}
		qApplicationApplication := qrApplicationApplication
		if qApplicationApplication != "" {

			if err := r.SetQueryParam("application.application", qApplicationApplication); err != nil {
				return err
			}
		}
	}

	if o.ApplicationProject != nil {

		// query param application.project
		var qrApplicationProject string

		if o.ApplicationProject != nil {
			qrApplicationProject = *o.ApplicationProject
		}
		qApplicationProject := qrApplicationProject
		if qApplicationProject != "" {

			if err := r.SetQueryParam("application.project", qApplicationProject); err != nil {
				return err
			}
		}
	}

	if o.ProjectProject != nil {

		// query param project.project
		var qrProjectProject string

		if o.ProjectProject != nil {
			qrProjectProject = *o.ProjectProject
		}
		qProjectProject := qrProjectProject
		if qProjectProject != "" {

			if err := r.SetQueryParam("project.project", qProjectProject); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
