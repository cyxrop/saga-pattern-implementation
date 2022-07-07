package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *ProductWarehouses) Update(ctx context.Context, pw model.ProductWarehouse) error {
	return r.update(ctx, r.pool, pw)
}

func (r *ProductWarehouses) UpdateWithTx(ctx context.Context, tx pgx.Tx, pw model.ProductWarehouse) error {
	return r.update(ctx, tx, pw)
}

func (r *ProductWarehouses) update(ctx context.Context, q pgxtype.Querier, pw model.ProductWarehouse) error {
	const query = `
		update product_warehouse
		set
			number = $3
		where
			product_id = $1 and
			warehouse_id = $2;
	`

	res, err := q.Exec(ctx, query, pw.ProductID, pw.WarehouseID, pw.Number)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}

func (r *ProductWarehouses) AddNumber(ctx context.Context, productID, warehouseID int64, diff int64) error {
	return r.addNumber(ctx, r.pool, productID, warehouseID, diff)
}

func (r *ProductWarehouses) AddNumberWithTx(ctx context.Context, tx pgx.Tx, productID, warehouseID int64, diff int64) error {
	return r.addNumber(ctx, tx, productID, warehouseID, diff)
}

func (r *ProductWarehouses) addNumber(ctx context.Context, q pgxtype.Querier, productID, warehouseID int64, diff int64) error {
	const query = `
		update product_warehouse
		set
			number = number + $3
		where
			product_id = $1 and
			warehouse_id = $2;
	`

	res, err := q.Exec(ctx, query, productID, warehouseID, diff)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
