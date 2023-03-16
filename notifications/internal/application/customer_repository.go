package application

import (
	"context"

	"github.com/50HJ/Intelli-Mall/notifications/internal/models"
)

type CustomerRepository interface {
	Find(ctx context.Context, customerID string) (*models.Customer, error)
}
