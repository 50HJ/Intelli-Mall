// Code generated by go-swagger; DO NOT EDIT.

package invoice

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/50HJ/Intelli-Mall/payments/paymentsclient/models"
)

// PayInvoiceReader is a Reader for the PayInvoice structure.
type PayInvoiceReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *PayInvoiceReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewPayInvoiceOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewPayInvoiceDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewPayInvoiceOK creates a PayInvoiceOK with default headers values
func NewPayInvoiceOK() *PayInvoiceOK {
	return &PayInvoiceOK{}
}

/* PayInvoiceOK describes a response with status code 200, with default header values.

A successful response.
*/
type PayInvoiceOK struct {
	Payload models.PaymentspbPayInvoiceResponse
}

// IsSuccess returns true when this pay invoice o k response has a 2xx status code
func (o *PayInvoiceOK) IsSuccess() bool {
	return true
}

// IsRedirect returns true when this pay invoice o k response has a 3xx status code
func (o *PayInvoiceOK) IsRedirect() bool {
	return false
}

// IsClientError returns true when this pay invoice o k response has a 4xx status code
func (o *PayInvoiceOK) IsClientError() bool {
	return false
}

// IsServerError returns true when this pay invoice o k response has a 5xx status code
func (o *PayInvoiceOK) IsServerError() bool {
	return false
}

// IsCode returns true when this pay invoice o k response a status code equal to that given
func (o *PayInvoiceOK) IsCode(code int) bool {
	return code == 200
}

func (o *PayInvoiceOK) Error() string {
	return fmt.Sprintf("[PUT /api/payments/invoices/{id}/pay][%d] payInvoiceOK  %+v", 200, o.Payload)
}

func (o *PayInvoiceOK) String() string {
	return fmt.Sprintf("[PUT /api/payments/invoices/{id}/pay][%d] payInvoiceOK  %+v", 200, o.Payload)
}

func (o *PayInvoiceOK) GetPayload() models.PaymentspbPayInvoiceResponse {
	return o.Payload
}

func (o *PayInvoiceOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewPayInvoiceDefault creates a PayInvoiceDefault with default headers values
func NewPayInvoiceDefault(code int) *PayInvoiceDefault {
	return &PayInvoiceDefault{
		_statusCode: code,
	}
}

/* PayInvoiceDefault describes a response with status code -1, with default header values.

An unexpected error response.
*/
type PayInvoiceDefault struct {
	_statusCode int

	Payload *models.RPCStatus
}

// Code gets the status code for the pay invoice default response
func (o *PayInvoiceDefault) Code() int {
	return o._statusCode
}

// IsSuccess returns true when this pay invoice default response has a 2xx status code
func (o *PayInvoiceDefault) IsSuccess() bool {
	return o._statusCode/100 == 2
}

// IsRedirect returns true when this pay invoice default response has a 3xx status code
func (o *PayInvoiceDefault) IsRedirect() bool {
	return o._statusCode/100 == 3
}

// IsClientError returns true when this pay invoice default response has a 4xx status code
func (o *PayInvoiceDefault) IsClientError() bool {
	return o._statusCode/100 == 4
}

// IsServerError returns true when this pay invoice default response has a 5xx status code
func (o *PayInvoiceDefault) IsServerError() bool {
	return o._statusCode/100 == 5
}

// IsCode returns true when this pay invoice default response a status code equal to that given
func (o *PayInvoiceDefault) IsCode(code int) bool {
	return o._statusCode == code
}

func (o *PayInvoiceDefault) Error() string {
	return fmt.Sprintf("[PUT /api/payments/invoices/{id}/pay][%d] payInvoice default  %+v", o._statusCode, o.Payload)
}

func (o *PayInvoiceDefault) String() string {
	return fmt.Sprintf("[PUT /api/payments/invoices/{id}/pay][%d] payInvoice default  %+v", o._statusCode, o.Payload)
}

func (o *PayInvoiceDefault) GetPayload() *models.RPCStatus {
	return o.Payload
}

func (o *PayInvoiceDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.RPCStatus)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
