package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Warehouses struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Warehouses {
	return &Warehouses{
		pool: p,
	}
}
