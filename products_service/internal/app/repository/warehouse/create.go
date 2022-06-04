package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
)

func (r *Warehouses) Create(ctx context.Context, w model.Warehouse) (int64, error) {
	const query = `
		insert into warehouses (
			name,
			address,
			created_at
		) VALUES (
			$1, $2, now()
		) returning id
	`

	var ID int64
	return ID, r.pool.QueryRow(ctx, query, w.Name, w.Address).Scan(&ID)
}
