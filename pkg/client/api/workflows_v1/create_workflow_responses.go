// Code generated by go-swagger; DO NOT EDIT.

package workflows_v1

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/puppetlabs/nebula/pkg/client/api/models"
)

// CreateWorkflowReader is a Reader for the CreateWorkflow structure.
type CreateWorkflowReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateWorkflowReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 201:
		result := NewCreateWorkflowCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewCreateWorkflowCreated creates a CreateWorkflowCreated with default headers values
func NewCreateWorkflowCreated() *CreateWorkflowCreated {
	return &CreateWorkflowCreated{}
}

/*CreateWorkflowCreated handles this case with default header values.

Newly created workflow
*/
type CreateWorkflowCreated struct {
	Payload *models.ShowWorkflow
}

func (o *CreateWorkflowCreated) Error() string {
	return fmt.Sprintf("[POST /api/workflows][%d] createWorkflowCreated  %+v", 201, o.Payload)
}

func (o *CreateWorkflowCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ShowWorkflow)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
