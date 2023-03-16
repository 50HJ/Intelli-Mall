package application

import (
	"context"

	"github.com/50HJ/Intelli-Mall/payments/internal/models"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *models.Payment) error
	Find(ctx context.Context, paymentID string) (*models.Payment, error)
}
