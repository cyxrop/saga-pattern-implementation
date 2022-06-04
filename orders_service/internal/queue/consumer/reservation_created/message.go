package consumer

type ReservationCreatedMessage struct {
	OrderID int64 `json:"order_id"`
}
