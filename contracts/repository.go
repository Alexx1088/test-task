package contracts

import (
	"context"

	"github.com/Alexx1088/test-task/domain"
)

type OrderRepository interface {
	Retrieve(ctx context.Context, orderID string) (*domain.Order, error)
	Update(ctx context.Context, order *domain.Order) error
	CreateAuditLog(ctx context.Context, orderID, action, value string) error
	RunInTx(ctx context.Context, fn func(repo OrderRepository) error) error
}
