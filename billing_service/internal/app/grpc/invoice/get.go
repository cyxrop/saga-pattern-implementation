package grpc

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/internal/app/model"
	"gitlab.ozon.dev/cyxrop/homework-3/billing_service/pkg/api"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

var statusMap = map[model.InvoiceStatus]api.InvoiceStatus{
	model.InvoiceStatusPendingPayment: api.InvoiceStatus_PendingPayment,
	model.InvoiceStatusPaid:           api.InvoiceStatus_Paid,
	model.InvoiceStatusFailed:         api.InvoiceStatus_Failed,
}

func (s *InvoiceServiceServer) Get(ctx context.Context, r *api.ID) (*api.InvoiceResponse, error) {
	invoice, err := s.service.GetInvoice(ctx, r.ID)
	if err != nil {
		return nil, err
	}

	status, err := mapStatus(invoice.Status)
	if err != nil {
		return nil, err
	}

	return &api.InvoiceResponse{
		ID:        invoice.ID,
		OrderID:   invoice.OrderID,
		Amount:    invoice.Amount,
		Status:    status,
		CreatedAt: tspb.New(invoice.CreatedAt),
	}, nil
}

func mapStatus(s model.InvoiceStatus) (api.InvoiceStatus, error) {
	status, ok := statusMap[s]
	if !ok {
		return api.InvoiceStatus_Failed, fmt.Errorf("unknown status: %d", s)
	}

	return status, nil
}
