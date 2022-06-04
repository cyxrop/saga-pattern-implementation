package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
)

func (r *Reservations) Create(ctx context.Context, rv model.Reservation) (int64, error) {
	return r.create(ctx, r.pool, rv)
}

func (r *Reservations) CreateWithTx(ctx context.Context, tx pgx.Tx, rv model.Reservation) (int64, error) {
	return r.create(ctx, tx, rv)
}

func (r *Reservations) create(ctx context.Context, q pgxtype.Querier, rv model.Reservation) (int64, error) {
	const query = `
		insert into reservations (
			product_id,
			warehouse_id,
			number,
			order_id,
			created_at
		) VALUES (
			$1, $2, $3, $4, now()
		) returning id
	`

	var ID int64
	return ID, q.QueryRow(ctx, query, rv.ProductID, rv.WarehouseID, rv.Number, rv.OrderID).Scan(&ID)
}
