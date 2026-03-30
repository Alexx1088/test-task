package bulk_complete

import (
	"context"
	"fmt"
	"github.com/Alexx1088/test-task/contracts"
	"strings"
	"sync"
)

type Interactor struct {
	repo contracts.OrderRepository
}

func NewInteractor(repo contracts.OrderRepository) *Interactor {
	return &Interactor{repo: repo}
}

func (uc *Interactor) Execute(ctx context.Context, orderIDs []string) error {
	const maxWorkers = 10

	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	var mu sync.Mutex
	var errs []string

	for _, id := range orderIDs {
		if ctx.Err() != nil {
			break
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(orderID string) {
			defer wg.Done()
			defer func() { <-sem }()

			order, err := uc.repo.Retrieve(ctx, orderID)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: retrieve failed: %v", orderID, err))
				mu.Unlock()
				return
			}

			if err := order.Complete(); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: complete failed: %v", orderID, err))
				mu.Unlock()
				return
			}

			if err := uc.repo.Update(ctx, order); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Sprintf("%s: update failed: %v", orderID, err))
				mu.Unlock()
				return
			}
		}(id)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("bulk complete failed: %s", strings.Join(errs, "; "))
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}
