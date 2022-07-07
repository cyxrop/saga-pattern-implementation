package consumer

type ReservationFailedMessage struct {
	OrderID int64 `json:"order_id"`
}
