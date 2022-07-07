package consumer

type InvoicePaidMessage struct {
	InvoiceID int64 `json:"invoice_id"`
	OrderID   int64 `json:"order_id"`
	Amount    int64 `json:"amount"`
}

type OrderPaidMessage struct {
	InvoiceID int64 `json:"invoice_id"`
	OrderID   int64 `json:"order_id"`
	Amount    int64 `json:"amount"`
}
