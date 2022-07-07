package model

import "time"

type Reservation struct {
	ID          int64
	ProductID   int64
	WarehouseID int64
	Number      int64
	OrderID     int64
	CreatedAt   time.Time
}
