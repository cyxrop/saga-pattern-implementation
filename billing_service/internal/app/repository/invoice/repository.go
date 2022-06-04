package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Invoices struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Invoices {
	return &Invoices{
		pool: p,
	}
}
