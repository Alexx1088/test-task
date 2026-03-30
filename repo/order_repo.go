package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Alexx1088/test-task/contracts"
	"github.com/Alexx1088/test-task/domain"
)

type OrderRepo struct {
	db *sql.DB
	tx *sql.Tx
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Retrieve(ctx context.Context, orderID string) (*domain.Order, error) {
	row := r.queryRowContext(ctx, `
		SELECT id, customer_id, status, total_cents
		FROM orders
		WHERE id = ?
	`, orderID)

	var order domain.Order
	if err := row.Scan(&order.ID, &order.CustomerID, &order.Status, &order.TotalCents); err != nil {
		return nil, fmt.Errorf("scan order: %w", err)
	}

	return &order, nil
}

func (r *OrderRepo) Update(ctx context.Context, order *domain.Order) error {
	_, err := r.execContext(ctx, `
		UPDATE orders
		SET customer_id = ?, status = ?, total_cents = ?
		WHERE id = ?
	`, order.CustomerID, order.Status, order.TotalCents, order.ID)
	if err != nil {
		return fmt.Errorf("update order: %w", err)
	}

	return nil
}

func (r *OrderRepo) CreateAuditLog(ctx context.Context, orderID, action, value string) error {
	_, err := r.execContext(ctx, `
		INSERT INTO audit_logs (order_id, action, value)
		VALUES (?, ?, ?)
	`, orderID, action, value)
	if err != nil {
		return fmt.Errorf("insert audit log: %w", err)
	}

	return nil
}

func (r *OrderRepo) RunInTx(ctx context.Context, fn func(repo contracts.OrderRepository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	txRepo := &OrderRepo{
		db: r.db,
		tx: tx,
	}

	if err := fn(txRepo); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func (r *OrderRepo) queryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if r.tx != nil {
		return r.tx.QueryRowContext(ctx, query, args...)
	}
	return r.db.QueryRowContext(ctx, query, args...)
}

func (r *OrderRepo) execContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if r.tx != nil {
		return r.tx.ExecContext(ctx, query, args...)
	}
	return r.db.ExecContext(ctx, query, args...)
}
