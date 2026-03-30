package domain

import (
	"errors"
)

type OrderStatus string

const (
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCompleted OrderStatus = "COMPLETED"
)

var (
	ErrOrderCannotBeCompleted = errors.New("order cannot be completed")
	ErrEmptyCustomerID        = errors.New("customer id cannot be empty")
)

type Order struct {
	ID         string
	CustomerID string
	Status     OrderStatus
	TotalCents int64
}

func (o *Order) Complete() error {
	if o.Status != StatusConfirmed {
		return ErrOrderCannotBeCompleted
	}

	o.Status = StatusCompleted
	return nil
}

func (o *Order) TransferTo(newCustomerID string) error {
	if newCustomerID == "" {
		return ErrEmptyCustomerID
	}

	o.CustomerID = newCustomerID
	return nil
}
