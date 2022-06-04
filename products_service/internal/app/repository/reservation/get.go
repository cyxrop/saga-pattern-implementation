package repository

import (
	"context"

	"github.com/jackc/pgtype/pgxtype"
	"github.com/jackc/pgx/v4"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
)

func (r *Reservations) GetOrderReservations(ctx context.Context, orderID int64) ([]model.Reservation, error) {
	return r.getOrderReservations(ctx, r.pool, orderID)
}

func (r *Reservations) GetOrderReservationsWithTx(ctx context.Context, tx pgx.Tx, orderID int64) ([]model.Reservation, error) {
	return r.getOrderReservations(ctx, tx, orderID)
}

func (r *Reservations) getOrderReservations(ctx context.Context, q pgxtype.Querier, orderID int64) ([]model.Reservation, error) {
	const query = `
		select
			id,
			product_id,
			warehouse_id,
			number,
			order_id,
			created_at
		from
			reservations
		where
			order_id = $1;
	`

	rows, err := q.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}

	var reservations []model.Reservation
	for rows.Next() {
		var res model.Reservation
		if err = rows.Scan(
			&res.ID,
			&res.ProductID,
			&res.WarehouseID,
			&res.Number,
			&res.OrderID,
			&res.CreatedAt,
		); err != nil {
			return nil, err
		}

		reservations = append(reservations, res)
	}

	return reservations, nil
}
