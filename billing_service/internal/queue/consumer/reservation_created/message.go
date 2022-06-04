package consumer

type ReservationCreatedMessage struct {
	OrderID int64 `json:"order_id"`
	Amount  int64 `json:"amount"`
}

type InvoiceIssuedMessage struct {
	InvoiceID int64 `json:"invoice_id"`
	OrderID   int64 `json:"order_id"`
	Amount    int64 `json:"amount"`
}
