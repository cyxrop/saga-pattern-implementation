package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/repository"
)

func (r *Invoices) Get(ctx context.Context, ID int64) (model.Invoice, error) {
	return r.get(ctx, r.pool, ID)
}

func (r *Invoices) GetWithTx(ctx context.Context, tx pgx.Tx, ID int64) (model.Invoice, error) {
	return r.get(ctx, tx, ID)
}

func (r *Invoices) get(ctx context.Context, q pgxtype.Querier, ID int64) (model.Invoice, error) {
	const query = `
		select
			id,
			order_id,
			amount,
			status,
			created_at
		from
			invoices
		where
			id = $1;
	`

	var i model.Invoice
	err := q.QueryRow(ctx, query, ID).Scan(
		&i.ID,
		&i.OrderID,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return i, repository.ErrNotFound
	}

	return i, err
}
