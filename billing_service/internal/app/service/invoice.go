package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/cache"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
	repository "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/repository/invoice"
)

type Invoices struct {
	ir    *repository.Invoices
	cache *cache.Cache
}

func New(ir *repository.Invoices, cache *cache.Cache) Invoices {
	return Invoices{
		ir:    ir,
		cache: cache,
	}
}

func (s *Invoices) CreateInvoice(ctx context.Context, orderID, amount int64) (int64, error) {
	invoice := model.Invoice{
		OrderID: orderID,
		Amount:  amount,
		Status:  model.InvoiceStatusPendingPayment,
	}

	ID, err := s.ir.Create(ctx, model.Invoice{
		OrderID: orderID,
		Amount:  amount,
		Status:  model.InvoiceStatusPendingPayment,
	})
	if err != nil {
		return 0, fmt.Errorf("create invoice: %w", err)
	}

	invoice.ID = ID

	if err := s.saveToCache(invoice); err != nil {
		log.Printf("failed save invoice %d to cache: %s", invoice.ID, err)
	}

	return ID, nil
}

func (s *Invoices) UpdateInvoiceStatus(ctx context.Context, ID int64, status model.InvoiceStatus) error {
	if err := s.cache.Delete(strconv.FormatInt(ID, 10)); err != nil {
		log.Printf("failed to delete cached invoice %d: %s", ID, err)
	}

	return s.ir.UpdateStatus(ctx, ID, status)
}

func (s *Invoices) GetInvoice(ctx context.Context, ID int64) (model.Invoice, error) {
	cached, err := s.cache.Get(strconv.FormatInt(ID, 10))
	if err == nil {
		var invoice model.Invoice
		if err = json.Unmarshal(cached, &invoice); err == nil {
			return invoice, nil
		}

		log.Printf("failed to unmarshal cached invoice %d: %s", ID, err)
	}

	invoice, err := s.ir.Get(ctx, ID)
	if err != nil {
		return model.Invoice{}, err
	}

	if err = s.saveToCache(invoice); err != nil {
		log.Printf("failed save invoice %d to cache: %s", invoice.ID, err)
	}

	return invoice, nil
}

func (s *Invoices) saveToCache(i model.Invoice) error {
	marshaled, err := json.Marshal(i)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return s.cache.Set(strconv.FormatInt(i.ID, 10), marshaled)
}
