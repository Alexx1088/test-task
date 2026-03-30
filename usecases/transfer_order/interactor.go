package transfer_order

import (
	"context"
	"fmt"
	"github.com/Alexx1088/test-task/contracts"
)

type Interactor struct {
	repo contracts.OrderRepository
}

func NewInteractor(repo contracts.OrderRepository) *Interactor {
	return &Interactor{repo: repo}
}

func (uc *Interactor) Execute(ctx context.Context, orderID, newCustomerID string) error {
	return uc.repo.RunInTx(ctx, func(repo contracts.OrderRepository) error {
		order, err := repo.Retrieve(ctx, orderID)
		if err != nil {
			return fmt.Errorf("retrieve order: %w", err)
		}

		if err := order.TransferTo(newCustomerID); err != nil {
			return fmt.Errorf("transfer order: %w", err)
		}

		if err := repo.Update(ctx, order); err != nil {
			return fmt.Errorf("update order: %w", err)
		}

		if err := repo.CreateAuditLog(ctx, orderID, "TRANSFERRED", newCustomerID); err != nil {
			return fmt.Errorf("create audit log: %w", err)
		}

		return nil
	})
}
