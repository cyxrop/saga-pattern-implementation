package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
)

func (r *OrderProducts) Create(ctx context.Context, op model.OrderProduct) error {
	return r.create(ctx, r.pool, op)
}

func (r *OrderProducts) CreateWithTx(ctx context.Context, tx pgx.Tx, op model.OrderProduct) error {
	return r.create(ctx, tx, op)
}

func (r *OrderProducts) create(ctx context.Context, q pgxtype.Querier, op model.OrderProduct) error {
	const query = `
		insert into order_products (
			order_id,
			product_id,
			number
		) VALUES (
			$1, $2, $3
		)
	`

	_, err := q.Exec(ctx, query, op.OrderID, op.ProductID, op.Number)
	return err
}
