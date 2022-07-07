package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Products struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Products {
	return &Products{
		pool: p,
	}
}
