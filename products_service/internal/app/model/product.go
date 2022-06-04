package model

import "time"

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       int64
	CreatedAt   time.Time
}

type ProductItem struct {
	ProductID int64
	Number    int64
}
