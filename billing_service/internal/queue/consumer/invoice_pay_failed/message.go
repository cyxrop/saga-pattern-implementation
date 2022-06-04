package consumer

type InvoicePayFailedMessage struct {
	InvoiceID int64 `json:"invoice_id"`
	OrderID   int64 `json:"order_id"`
	Amount    int64 `json:"amount"`
}

type OrderPayFailedMessage struct {
	OrderID int64 `json:"order_id"`
}
