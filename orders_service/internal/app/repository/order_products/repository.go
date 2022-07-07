package repository

import "github.com/jackc/pgx/v4/pgxpool"

type OrderProducts struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *OrderProducts {
	return &OrderProducts{
		pool: p,
	}
}
