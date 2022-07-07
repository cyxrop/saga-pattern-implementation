package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *Warehouses) Update(ctx context.Context, w model.Warehouse) error {
	const query = `
		update warehouses
		set
			name = $2,
			address = $3
		where id = $1;
	`

	res, err := r.pool.Exec(ctx, query, w.ID, w.Name, w.Address)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
