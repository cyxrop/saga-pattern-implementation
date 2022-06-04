package repository

import "github.com/jackc/pgx/v4/pgxpool"

type Reservations struct {
	pool *pgxpool.Pool
}

func New(p *pgxpool.Pool) *Reservations {
	return &Reservations{
		pool: p,
	}
}
