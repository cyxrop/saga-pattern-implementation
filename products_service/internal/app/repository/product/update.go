package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *Products) Update(ctx context.Context, p model.Product) error {
	const query = `
		update products
		set
			name = $2,
			description = $3,
			price = $4
		where id = $1;
	`

	res, err := r.pool.Exec(ctx, query, p.ID, p.Name, p.Description, p.Price)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
