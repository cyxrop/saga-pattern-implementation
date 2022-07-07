package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/repository"
)

func (r *Orders) Update(ctx context.Context, o model.Order) error {
	return r.update(ctx, r.pool, o)
}

func (r *Orders) UpdateWithTx(ctx context.Context, tx pgx.Tx, o model.Order) error {
	return r.update(ctx, tx, o)
}

func (r *Orders) update(ctx context.Context, q pgxtype.Querier, o model.Order) error {
	const query = `
		update orders
		set
			warehouse_id = $2,
			status = $3,
			status_description = $4
		where id = $1;
	`

	res, err := q.Exec(ctx, query, o.WarehouseID, o.Status, o.StatusDescription)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *Orders) UpdateStatus(ctx context.Context, ID int64, s model.OrderStatus, desc string) error {
	const query = `
		update orders
		set
			status = $2,
			status_description = $3
		where id = $1;
	`

	res, err := r.pool.Exec(ctx, query, ID, s, desc)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
