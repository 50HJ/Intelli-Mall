// Code generated by go-swagger; DO NOT EDIT.

package shopping_list

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/50HJ/Intelli-Mall/depot/depotclient/models"
)

// CcompleteShoppingListReader is a Reader for the CcompleteShoppingList structure.
type CcompleteShoppingListReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CcompleteShoppingListReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCcompleteShoppingListOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewCcompleteShoppingListDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCcompleteShoppingListOK creates a CcompleteShoppingListOK with default headers values
func NewCcompleteShoppingListOK() *CcompleteShoppingListOK {
	return &CcompleteShoppingListOK{}
}

/* CcompleteShoppingListOK describes a response with status code 200, with default header values.

A successful response.
*/
type CcompleteShoppingListOK struct {
	Payload models.DepotpbCompleteShoppingListResponse
}

// IsSuccess returns true when this ccomplete shopping list o k response has a 2xx status code
func (o *CcompleteShoppingListOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this ccomplete shopping list o k response has a 3xx status code
func (o *CcompleteShoppingListOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this ccomplete shopping list o k response has a 4xx status code
func (o *CcompleteShoppingListOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this ccomplete shopping list o k response has a 5xx status code
func (o *CcompleteShoppingListOK) IsServerError() bool {
	return false
}

// IsCode returns true when this ccomplete shopping list o k response a status code equal to that given
func (o *CcompleteShoppingListOK) IsCode(code int) bool {
	return code == 200
}

func (o *CcompleteShoppingListOK) Error() string {
	return fmt.Sprintf("[PUT /api/depot/shopping/{id}/complete][%d] ccompleteShoppingListOK  %+v", 200, o.Payload)
}

func (o *CcompleteShoppingListOK) String() string {
	return fmt.Sprintf("[PUT /api/depot/shopping/{id}/complete][%d] ccompleteShoppingListOK  %+v", 200, o.Payload)
}

func (o *CcompleteShoppingListOK) GetPayload() models.DepotpbCompleteShoppingListResponse {
	return o.Payload
}

func (o *CcompleteShoppingListOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCcompleteShoppingListDefault creates a CcompleteShoppingListDefault with default headers values
func NewCcompleteShoppingListDefault(code int) *CcompleteShoppingListDefault {
	return &CcompleteShoppingListDefault{
		_statusCode: code,
	}
}

/* CcompleteShoppingListDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type CcompleteShoppingListDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// Code gets the status code for the ccomplete shopping list default response
func (o *CcompleteShoppingListDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this ccomplete shopping list default response has a 2xx status code
func (o *CcompleteShoppingListDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this ccomplete shopping list default response has a 3xx status code
func (o *CcompleteShoppingListDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this ccomplete shopping list default response has a 4xx status code
func (o *CcompleteShoppingListDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this ccomplete shopping list default response has a 5xx status code
func (o *CcompleteShoppingListDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this ccomplete shopping list default response a status code equal to that given
func (o *CcompleteShoppingListDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *CcompleteShoppingListDefault) Error() string {
	return fmt.Sprintf("[PUT /api/depot/shopping/{id}/complete][%d] ccompleteShoppingList default  %+v", o._statusCode, o.Payload)
}

func (o *CcompleteShoppingListDefault) String() string {
	return fmt.Sprintf("[PUT /api/depot/shopping/{id}/complete][%d] ccompleteShoppingList default  %+v", o._statusCode, o.Payload)
}

func (o *CcompleteShoppingListDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *CcompleteShoppingListDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
