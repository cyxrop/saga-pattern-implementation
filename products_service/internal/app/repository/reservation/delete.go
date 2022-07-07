package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *Reservations) Delete(ctx context.Context, ID int64) error {
	const query = `
		delete from reservations
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

func (r *Reservations) DeleteByIDs(ctx context.Context, IDs []int64) error {
	return r.deleteByIDs(ctx, r.pool, IDs)
}

func (r *Reservations) DeleteByIDsWithTx(ctx context.Context, tx pgx.Tx, IDs []int64) error {
	return r.deleteByIDs(ctx, tx, IDs)
}

func (r *Reservations) deleteByIDs(ctx context.Context, q pgxtype.Querier, IDs []int64) error {
	const query = `
		delete from reservations
		where id = any($1);
	`

	res, err := q.Exec(ctx, query, pq.Array(IDs))
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
