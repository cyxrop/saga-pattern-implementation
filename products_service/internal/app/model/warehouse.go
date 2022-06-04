package model

import "time"

type Warehouse struct {
	ID        int64
	Name      string
	Address   string
	CreatedAt time.Time
}
