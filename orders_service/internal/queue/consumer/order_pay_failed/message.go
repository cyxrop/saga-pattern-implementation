package consumer

type OrderPayFailedMessage struct {
	OrderID int64 `json:"order_id"`
}
