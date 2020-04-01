// Code generated by go-swagger; DO NOT EDIT.

package notifications

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// New creates a new notifications API client.
func New(transport runtime.ClientTransport, formats strfmt.Registry) ClientService {
	return &Client{transport: transport, formats: formats}
}

/*
Client for notifications API
*/
type Client struct {
	transport runtime.ClientTransport
	formats   strfmt.Registry
}

// ClientService is the interface for Client methods
type ClientService interface {
	GetNotifications(params *GetNotificationsParams, authInfo runtime.ClientAuthInfoWriter) (*GetNotificationsOK, error)

	PostAllNotificationRead(params *PostAllNotificationReadParams, authInfo runtime.ClientAuthInfoWriter) (*PostAllNotificationReadOK, error)

	PostNotificationRead(params *PostNotificationReadParams, authInfo runtime.ClientAuthInfoWriter) (*PostNotificationReadOK, error)

	SetTransport(transport runtime.ClientTransport)
}

/*
  GetNotifications gets a list of notifications
*/
func (a *Client) GetNotifications(params *GetNotificationsParams, authInfo runtime.ClientAuthInfoWriter) (*GetNotificationsOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewGetNotificationsParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "getNotifications",
		Method:             "GET",
		PathPattern:        "/api/notifications",
		ProducesMediaTypes: []string{"application/vnd.puppet.nebula.v20200131+json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &GetNotificationsReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*GetNotificationsOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*GetNotificationsDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  PostAllNotificationRead marks multiple notification as read
*/
func (a *Client) PostAllNotificationRead(params *PostAllNotificationReadParams, authInfo runtime.ClientAuthInfoWriter) (*PostAllNotificationReadOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostAllNotificationReadParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "postAllNotificationRead",
		Method:             "POST",
		PathPattern:        "/api/notifications/read",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/vnd.puppet.nebula.v20200131+json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostAllNotificationReadReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostAllNotificationReadOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PostAllNotificationReadDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

/*
  PostNotificationRead marks notification as read
*/
func (a *Client) PostNotificationRead(params *PostNotificationReadParams, authInfo runtime.ClientAuthInfoWriter) (*PostNotificationReadOK, error) {
	// TODO: Validate the params before sending
	if params == nil {
		params = NewPostNotificationReadParams()
	}

	result, err := a.transport.Submit(&runtime.ClientOperation{
		ID:                 "postNotificationRead",
		Method:             "POST",
		PathPattern:        "/api/notifications/read/{notificationId}",
		ProducesMediaTypes: []string{"application/json"},
		ConsumesMediaTypes: []string{"application/json"},
		Schemes:            []string{"https"},
		Params:             params,
		Reader:             &PostNotificationReadReader{formats: a.formats},
		AuthInfo:           authInfo,
		Context:            params.Context,
		Client:             params.HTTPClient,
	})
	if err != nil {
		return nil, err
	}
	success, ok := result.(*PostNotificationReadOK)
	if ok {
		return success, nil
	}
	// unexpected success response
	unexpectedSuccess := result.(*PostNotificationReadDefault)
	return nil, runtime.NewAPIError("unexpected success response: content available as default response in error", unexpectedSuccess, unexpectedSuccess.Code())
}

// SetTransport changes the transport on the client
func (a *Client) SetTransport(transport runtime.ClientTransport) {
	a.transport = transport
}
