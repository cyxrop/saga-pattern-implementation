package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
)

func (r *Orders) Create(ctx context.Context, o model.Order) (int64, error) {
	return r.create(ctx, r.pool, o)
}

func (r *Orders) CreateWithTx(ctx context.Context, tx pgx.Tx, o model.Order) (int64, error) {
	return r.create(ctx, tx, o)
}

func (r *Orders) create(ctx context.Context, q pgxtype.Querier, o model.Order) (int64, error) {
	const query = `
		insert into orders (
			warehouse_id,
			status,
			status_description,
			created_at
		) VALUES (
			$1, $2, $3, now()
		) returning id
	`

	var ID int64
	return ID, q.QueryRow(ctx, query, o.WarehouseID, o.Status, o.StatusDescription).Scan(&ID)
}
