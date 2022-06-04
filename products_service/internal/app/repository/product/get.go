package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"github.com/lib/pq"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
)

func (r *Products) Get(ctx context.Context, ID int64) (model.Product, error) {
	return r.get(ctx, r.pool, ID)
}

func (r *Products) GetWithTx(ctx context.Context, tx pgx.Tx, ID int64) (model.Product, error) {
	return r.get(ctx, tx, ID)
}

func (r *Products) get(ctx context.Context, q pgxtype.Querier, ID int64) (model.Product, error) {
	const query = `
		select
			id,
			name,
			description,
			price,
			created_at
		from
			products
		where
			id = $1;
	`

	var p model.Product
	err := q.QueryRow(ctx, query, ID).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, repository.ErrNotFound
	}

	return p, err
}

func (r *Products) GetListByIDs(ctx context.Context, IDs []int64) ([]model.Product, error) {
	return r.getListByIDs(ctx, r.pool, IDs)
}

func (r *Products) GetListByIDsWithTx(ctx context.Context, tx pgx.Tx, IDs []int64) ([]model.Product, error) {
	return r.getListByIDs(ctx, tx, IDs)
}

func (r *Products) getListByIDs(ctx context.Context, q pgxtype.Querier, IDs []int64) ([]model.Product, error) {
	const query = `
		select
			id,
			name,
			description,
			price,
			created_at
		from
			products
		where
			id = any($1);
	`

	rows, err := q.Query(ctx, query, pq.Array(IDs))
	if err != nil {
		return nil, err
	}

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}
