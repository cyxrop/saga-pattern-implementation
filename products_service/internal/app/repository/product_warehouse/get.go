package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *ProductWarehouses) Get(ctx context.Context, ProductID, WarehouseID int64) (model.ProductWarehouse, error) {
	return r.get(ctx, r.pool, ProductID, WarehouseID)
}

func (r *ProductWarehouses) GetWithTx(ctx context.Context, tx pgx.Tx, ProductID, WarehouseID int64) (model.ProductWarehouse, error) {
	return r.get(ctx, tx, ProductID, WarehouseID)
}

func (r *ProductWarehouses) get(ctx context.Context, q pgxtype.Querier, ProductID, WarehouseID int64) (model.ProductWarehouse, error) {
	const query = `
		select
			product_id,
			warehouse_id,
			number
		from product_warehouse
		where
			product_id = $1 and
			warehouse_id = $2;
	`

	var pw model.ProductWarehouse
	err := q.QueryRow(ctx, query, ProductID, WarehouseID).Scan(
		&pw.ProductID,
		&pw.WarehouseID,
		&pw.Number,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return pw, repository.ErrNotFound
	}

	return pw, err
}
