package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Orders struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Orders {
	return &Orders{
		pool: p,
	}
}
