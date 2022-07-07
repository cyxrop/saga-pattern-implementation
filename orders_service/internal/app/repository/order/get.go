package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/repository"
)

func (r *Orders) Get(ctx context.Context, ID int64) (model.Order, error) {
	return r.get(ctx, r.pool, ID)
}

func (r *Orders) GetWithTx(ctx context.Context, tx pgx.Tx, ID int64) (model.Order, error) {
	return r.get(ctx, tx, ID)
}

func (r *Orders) get(ctx context.Context, q pgxtype.Querier, ID int64) (model.Order, error) {
	const query = `
		select
			id,
			warehouse_id,
			status,
			status_description,
			created_at
		from
			orders
		where
			id = $1;
	`

	var o model.Order
	err := q.QueryRow(ctx, query, ID).Scan(
		&o.ID,
		&o.WarehouseID,
		&o.Status,
		&o.StatusDescription,
		&o.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return o, repository.ErrNotFound
	}

	return o, err
}
