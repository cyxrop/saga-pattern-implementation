package service

import (
	"context"

	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
	repository "gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/repository/invoice"
)

type Invoices struct {
	ir *repository.Invoices
}

func New(ir *repository.Invoices) Invoices {
	return Invoices{ir: ir}
}

func (s Invoices) CreateInvoice(ctx context.Context, orderID, amount int64) (int64, error) {
	return s.ir.Create(ctx, model.Invoice{
		OrderID: orderID,
		Amount:  amount,
		Status:  model.InvoiceStatusPendingPayment,
	})
}

func (s Invoices) UpdateInvoiceStatus(ctx context.Context, ID int64, status model.InvoiceStatus) error {
	return s.ir.UpdateStatus(ctx, ID, status)
}
