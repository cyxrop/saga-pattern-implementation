package consumer

type OrderPaidMessage struct {
	OrderID int64 `json:"order_id"`
}
