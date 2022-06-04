package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
)

func (r *ProductWarehouses) Create(ctx context.Context, pw model.ProductWarehouse) error {
	const query = `
		insert into product_warehouse (
			product_id,
			warehouse_id,
			number
		) VALUES (
			$1, $2, $3
		)
	`

	_, err := r.pool.Exec(ctx, query, pw.ProductID, pw.WarehouseID, pw.Number)
	return err
}
