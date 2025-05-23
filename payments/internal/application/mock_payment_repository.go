// Code generated by mockery v2.14.0. DO NOT EDIT.

package application

import (
	context "context"
	models "github.com/50HJ/Intelli-Mall/payments/internal/models"

	mock "github.com/stretchr/testify/mock"
)

// MockPaymentRepository is an autogenerated mock type for the PaymentRepository type
type MockPaymentRepository struct {
	mock.Mock
}

// Find provides a mock function with given fields: ctx, paymentID
func (_m *MockPaymentRepository) Find(ctx context.Context, paymentID string) (*models.Payment, error) {
	ret := _m.Called(ctx, paymentID)

	var r0 *models.Payment
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Payment); ok {
		r0 = rf(ctx, paymentID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Payment)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, paymentID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Save provides a mock function with given fields: ctx, payment
func (_m *MockPaymentRepository) Save(ctx context.Context, payment *models.Payment) error {
	ret := _m.Called(ctx, payment)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Payment) error); ok {
		r0 = rf(ctx, payment)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewMockPaymentRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockPaymentRepository creates a new instance of MockPaymentRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockPaymentRepository(t mockConstructorTestingTNewMockPaymentRepository) *MockPaymentRepository {
	mock := &MockPaymentRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
