package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/repository"
)

func (r *Orders) Delete(ctx context.Context, ID int64) error {
	const query = `
		delete from orders
		where id = $1;
	`

	res, err := r.pool.Exec(ctx, query, ID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
