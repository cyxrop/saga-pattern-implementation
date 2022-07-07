package consumer

type OrderCreatedMessage struct {
	OrderID     int64                      `json:"order_id"`
	WarehouseID int64                      `json:"warehouse_id"`
	Products    []OrderCreatedProductsItem `json:"products"`
}

type OrderCreatedProductsItem struct {
	ProductID int64 `json:"product_id"`
	Number    int64 `json:"number"`
}

type ReservationCreatedMessage struct {
	OrderID        int64   `json:"order_id"`
	ReservationIDs []int64 `json:"reservations"`
	Amount         int64   `json:"amount"`
}
