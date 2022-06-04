package consumer

type ProductItem struct {
	ProductID int64 `json:"product_id"`
	Number    int64 `json:"number"`
}

type CreateOrderMessage struct {
	WarehouseID int64         `json:"warehouse_id"`
	Products    []ProductItem `json:"products"`
}

type OrderCreatedMessage struct {
	OrderID     int64         `json:"order_id"`
	WarehouseID int64         `json:"warehouse_id"`
	Products    []ProductItem `json:"products"`
}
