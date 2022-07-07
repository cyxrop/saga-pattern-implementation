package service

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/repository"
	oRepository "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/repository/order"
	opRepository "gitlab.ozon.dev/cyxrop/homework-3/order_service/internal/app/repository/order_products"
)

type Orders struct {
	or  *oRepository.Orders
	opr *opRepository.OrderProducts
	r   *repository.Repository
}

func NewOrders(or *oRepository.Orders, opr *opRepository.OrderProducts, r *repository.Repository) Orders {
	return Orders{
		or:  or,
		opr: opr,
		r:   r,
	}
}

func (s Orders) CreateOrder(ctx context.Context, order model.Order, products []model.OrderProduct) (ID int64, err error) {
	tx, err := s.r.BeginTx(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}

		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	order.Status = model.OrderStatusPendingReservation
	ID, err = s.or.CreateWithTx(ctx, tx, order)
	if err != nil {
		return 0, fmt.Errorf("create order: %s", err)
	}

	for _, p := range products {
		p.OrderID = ID
		if err = s.opr.CreateWithTx(ctx, tx, p); err != nil {
			return 0, fmt.Errorf("create order product %d: %s", p.ProductID, err)
		}
	}

	return ID, nil
}

func (s Orders) UpdateOrderStatus(ctx context.Context, ID int64, status model.OrderStatus, desc string) error {
	return s.or.UpdateStatus(ctx, ID, status, desc)
}
