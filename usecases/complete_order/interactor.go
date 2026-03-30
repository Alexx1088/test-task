package complete_order

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

func (uc *Interactor) Execute(ctx context.Context, orderID string) error {
	order, err := uc.repo.Retrieve(ctx, orderID)
	if err != nil {
		return fmt.Errorf("retrieve order: %w", err)
	}

	if err := order.Complete(); err != nil {
		return fmt.Errorf("complete order: %w", err)
	}

	if err := uc.repo.Update(ctx, order); err != nil {
		return fmt.Errorf("persist order: %w", err)
	}

	return nil
}
