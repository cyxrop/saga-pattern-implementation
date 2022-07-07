package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
)

func (r *Products) Create(ctx context.Context, p model.Product) (int64, error) {
	const query = `
		insert into products (
			name,
			description,
			price,
			created_at
		) VALUES (
			$1, $2, $3, now()
		) returning id
	`

	var ID int64
	return ID, r.pool.QueryRow(ctx, query, p.Name, p.Description, p.Price).Scan(&ID)
}
