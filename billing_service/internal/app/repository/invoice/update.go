package repository

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/repository"
)

func (r *Invoices) UpdateStatus(ctx context.Context, ID int64, s model.InvoiceStatus) error {
	const query = `
		update invoices
		set
			status = $2
		where id = $1;
	`

	res, err := r.pool.Exec(ctx, query, ID, s)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return repository.ErrNotFound
	}

	return nil
}
