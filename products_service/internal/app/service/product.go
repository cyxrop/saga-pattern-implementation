package service

import (
	"context"
	"errors"
	"fmt"

	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository"
	pRepository "gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository/product"
	pwRepository "gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository/product_warehouse"
	rRepository "gitlab.ozon.dev/cyxrop/homework-3/products_service/internal/app/repository/reservation"
)

type Products struct {
	pwr *pwRepository.ProductWarehouses
	rr  *rRepository.Reservations
	pr  *pRepository.Products
	r   *repository.Repository
}

func NewProductService(
	pwr *pwRepository.ProductWarehouses,
	rr *rRepository.Reservations,
	pr *pRepository.Products,
	r *repository.Repository,
) Products {
	return Products{
		pwr: pwr,
		rr:  rr,
		pr:  pr,
		r:   r,
	}
}

func (s Products) ReserveOrderProducts(
	ctx context.Context,
	orderID, warehouseID int64,
	productReservations []model.ProductItem,
) (reservations []int64, err error) {
	tx, err := s.r.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
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

	for _, pr := range productReservations {
		if pr.Number <= 0 {
			continue
		}

		pw, err := s.pwr.GetWithTx(ctx, tx, pr.ProductID, warehouseID)
		if err != nil {
			return nil, fmt.Errorf("get product-warehouse %d-%d: %w", pr.ProductID, warehouseID, err)
		}

		newNumber := pw.Number - pr.Number
		if newNumber < 0 {
			return nil, fmt.Errorf("not enough products %d in warehouse %d", pr.ProductID, warehouseID)
		}

		pw.Number = newNumber
		if err = s.pwr.UpdateWithTx(ctx, tx, pw); err != nil {
			return nil, fmt.Errorf("cannot update product-warehouse %d-%d: %w", pr.ProductID, warehouseID, err)
		}

		reservationID, err := s.rr.CreateWithTx(ctx, tx, model.Reservation{
			ProductID:   pw.ProductID,
			WarehouseID: pw.WarehouseID,
			Number:      pr.Number,
			OrderID:     orderID,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot reserve product %d: %w", pw.ProductID, err)
		}

		reservations = append(reservations, reservationID)
	}

	return reservations, nil
}

func (s Products) CalculateAmount(ctx context.Context, pItems []model.ProductItem) (int64, error) {
	IDs := make([]int64, len(pItems))
	productNumbers := make(map[int64]int64)
	for i, pi := range pItems {
		IDs[i] = pi.ProductID
		productNumbers[pi.ProductID] = pi.Number
	}

	products, err := s.pr.GetListByIDs(ctx, IDs)
	if err != nil {
		return 0, fmt.Errorf("get products: %w", err)
	}

	if len(products) != len(pItems) {
		return 0, errors.New("found unknown products")
	}

	var amount int64
	for _, p := range products {
		number, ok := productNumbers[p.ID]
		if !ok {
			return 0, fmt.Errorf("get product number %d", p.ID)
		}

		amount = amount + p.Price*number
	}

	return amount, nil
}

func (s Products) CancelOrderReservation(ctx context.Context, orderID int64) (err error) {
	tx, err := s.r.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
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

	reservations, err := s.rr.GetOrderReservationsWithTx(ctx, tx, orderID)
	if err != nil {
		return fmt.Errorf("get order reservations: %s", err)
	}

	IDs := make([]int64, len(reservations))
	for i, r := range reservations {
		IDs[i] = r.ID
		if err = s.pwr.AddNumberWithTx(ctx, tx, r.ProductID, r.WarehouseID, r.Number); err != nil {
			return fmt.Errorf("add %d products %d to warehouse %d: %w", r.Number, r.ProductID, r.WarehouseID, err)
		}
	}

	if err = s.rr.DeleteByIDsWithTx(ctx, tx, IDs); err != nil {
		return fmt.Errorf("delete order reservations: %s", err)
	}

	return nil
}
