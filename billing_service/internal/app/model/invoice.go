package model

import "time"

type InvoiceStatus int8

const (
	InvoiceStatusPendingPayment InvoiceStatus = iota + 1
	InvoiceStatusPaid
	InvoiceStatusFailed
)

type Invoice struct {
	ID        int64
	OrderID   int64
	Amount    int64
	Status    InvoiceStatus
	CreatedAt time.Time
}
