package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
)

func (r *Invoices) Create(ctx context.Context, i model.Invoice) (int64, error) {
	const query = `
		insert into invoices (
			order_id,
			amount,
			status,
			created_at
		) VALUES (
			$1, $2, $3, now()
		) returning id
	`

	var ID int64
	return ID, r.pool.QueryRow(ctx, query, i.OrderID, i.Amount, i.Status).Scan(&ID)
}
