package model

import "time"

type OrderStatus int8

const (
	OrderStatusPendingReservation OrderStatus = iota + 1
	OrderStatusPendingPayment
	OrderStatusFailed
	OrderStatusPaid
)

type Order struct {
	ID                int64
	WarehouseID       int64
	Status            OrderStatus
	StatusDescription string
	CreatedAt         time.Time
}
