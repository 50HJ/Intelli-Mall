// Code generated by go-swagger; DO NOT EDIT.

package store

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/50HJ/Intelli-Mall/stores/storesclient/models"
)

// GetStoresReader is a Reader for the GetStores structure.
type GetStoresReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetStoresReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetStoresOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetStoresDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetStoresOK creates a GetStoresOK with default headers values
func NewGetStoresOK() *GetStoresOK {
	return &GetStoresOK{}
}

/* GetStoresOK describes a response with status code 200, with default header values.

A successful response.
*/
type GetStoresOK struct {
	Payload *models.StorespbGetStoresResponse
}

// IsSuccess returns true when this get stores o k response has a 2xx status code
func (o *GetStoresOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this get stores o k response has a 3xx status code
func (o *GetStoresOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this get stores o k response has a 4xx status code
func (o *GetStoresOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this get stores o k response has a 5xx status code
func (o *GetStoresOK) IsServerError() bool {
	return false
}

// IsCode returns true when this get stores o k response a status code equal to that given
func (o *GetStoresOK) IsCode(code int) bool {
	return code == 200
}

func (o *GetStoresOK) Error() string {
	return fmt.Sprintf("[GET /api/stores][%d] getStoresOK  %+v", 200, o.Payload)
}

func (o *GetStoresOK) String() string {
	return fmt.Sprintf("[GET /api/stores][%d] getStoresOK  %+v", 200, o.Payload)
}

func (o *GetStoresOK) GetPayload() *models.StorespbGetStoresResponse {
	return o.Payload
}

func (o *GetStoresOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.StorespbGetStoresResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetStoresDefault creates a GetStoresDefault with default headers values
func NewGetStoresDefault(code int) *GetStoresDefault {
	return &GetStoresDefault{
		_statusCode: code,
	}
}

/* GetStoresDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type GetStoresDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// Code gets the status code for the get stores default response
func (o *GetStoresDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this get stores default response has a 2xx status code
func (o *GetStoresDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this get stores default response has a 3xx status code
func (o *GetStoresDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this get stores default response has a 4xx status code
func (o *GetStoresDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this get stores default response has a 5xx status code
func (o *GetStoresDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this get stores default response a status code equal to that given
func (o *GetStoresDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *GetStoresDefault) Error() string {
	return fmt.Sprintf("[GET /api/stores][%d] getStores default  %+v", o._statusCode, o.Payload)
}

func (o *GetStoresDefault) String() string {
	return fmt.Sprintf("[GET /api/stores][%d] getStores default  %+v", o._statusCode, o.Payload)
}

func (o *GetStoresDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *GetStoresDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
