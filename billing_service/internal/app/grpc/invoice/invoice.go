package grpc

import (
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/service"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/pkg/api"
)

type InvoiceServiceServer struct {
	api.UnimplementedInvoiceServiceServer

	service service.Invoices
}

func NewInvoiceServiceServer(service service.Invoices) *InvoiceServiceServer {
	return &InvoiceServiceServer{service: service}
}
