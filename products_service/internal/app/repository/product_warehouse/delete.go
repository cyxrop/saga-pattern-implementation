package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *ProductWarehouses) Delete(ctx context.Context, ProductID, WarehouseID int64) error {
	const query = `
		delete from product_warehouse
		where
			product_id = $1 and
			warehouse_id = $2;
	`

	res, err := r.pool.Exec(ctx, query, ProductID, WarehouseID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
