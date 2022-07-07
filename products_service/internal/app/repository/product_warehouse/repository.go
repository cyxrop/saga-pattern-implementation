package repository

import "github.com/jackc/pgx/v4/pgxpool"

type ProductWarehouses struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *ProductWarehouses {
	return &ProductWarehouses{
		pool: p,
	}
}
