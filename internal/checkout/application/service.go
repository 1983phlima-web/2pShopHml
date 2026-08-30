package application

import (
	"context"
	"fmt"

	catalogDomain "github.com/2pshop/2pshop/internal/catalog/domain"
	catalogPorts "github.com/2pshop/2pshop/internal/catalog/ports"
	inventoryApp "github.com/2pshop/2pshop/internal/inventory/application"
	inventoryDomain "github.com/2pshop/2pshop/internal/inventory/domain"
	ordersApp "github.com/2pshop/2pshop/internal/orders/application"
	orderDomain "github.com/2pshop/2pshop/internal/orders/domain"
	paymentsApp "github.com/2pshop/2pshop/internal/payments/application"
	paymentDomain "github.com/2pshop/2pshop/internal/payments/domain"
	"github.com/2pshop/2pshop/pkg/errors"
	"github.com/2pshop/2pshop/pkg/idempotency"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	catalogRepo      catalogPorts.Repository
	inventoryService *inventoryApp.Service
	orderService     *ordersApp.Service
	paymentService   *paymentsApp.Service
	idempotency      *idempotency.Manager
	tracer           trace.Tracer
}

func NewService(
	catalogRepo catalogPorts.Repository,
	inventoryService *inventoryApp.Service,
	orderService *ordersApp.Service,
	paymentService *paymentsApp.Service,
	idempotencyMgr *idempotency.Manager,
	tracer trace.Tracer,
) *Service {
	return &Service{
		catalogRepo:      catalogRepo,
		inventoryService: inventoryService,
		orderService:     orderService,
		paymentService:   paymentService,
		idempotency:      idempotencyMgr,
		tracer:           tracer,
	}
}

type CheckoutRequest struct {
	TenantID       string                      `json:"tenant_id"`
	CustomerID     string                      `json:"customer_id"`
	Items          []CheckoutItem              `json:"items"`
	PaymentMethod  paymentDomain.PaymentMethod `json:"payment_method"`
	IdempotencyKey string                      `json:"idempotency_key"`
}

type CheckoutItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type CheckoutResult struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

// idempotencyFail is a nil-safe wrapper: idempotency tracking is optional
// (only enabled when REDIS_URL is configured), so failures must never
// panic when it's disabled.
func (s *Service) idempotencyFail(ctx context.Context, tenantID, key string) {
	if s.idempotency == nil || key == "" {
		return
	}
	_ = s.idempotency.Fail(ctx, tenantID, key)
}

func (s *Service) Checkout(ctx context.Context, req CheckoutRequest) (*CheckoutResult, error) {
	ctx, span := s.tracer.Start(ctx, "checkout.execute",
		trace.WithAttributes(
			attribute.String("tenant.id", req.TenantID),
			attribute.String("customer.id", req.CustomerID),
		),
	)
	defer span.End()

	// Idempotency check (only when Redis-backed idempotency is configured).
	if s.idempotency != nil && req.IdempotencyKey != "" {
		existing, err := s.idempotency.Check(ctx, req.TenantID, req.IdempotencyKey, "checkout", req)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return &CheckoutResult{OrderID: "", Status: string(existing.Status)}, nil
		}
	}

	// 1. Validate
	ctx, validateSpan := s.tracer.Start(ctx, "checkout.validate")
	var orderItems []orderDomain.OrderItem
	var total int64
	for _, item := range req.Items {
		product, err := s.catalogRepo.GetByID(ctx, req.TenantID, item.ProductID)
		if err != nil {
			validateSpan.RecordError(err)
			validateSpan.End()
			s.idempotencyFail(ctx, req.TenantID, req.IdempotencyKey)
			return nil, errors.Wrap(errors.ErrNotFound, "product not found", err)
		}
		if product.State != catalogDomain.StateActive {
			validateSpan.End()
			s.idempotencyFail(ctx, req.TenantID, req.IdempotencyKey)
			return nil, errors.New(errors.ErrInvalidInput).WithDetail("product", "not active")
		}
		orderItems = append(orderItems, orderDomain.OrderItem{
			ProductID: product.ID,
			Name:      product.Name,
			Quantity:  item.Quantity,
			UnitPrice: product.Price,
		})
		total += product.Price * int64(item.Quantity)
	}
	validateSpan.End()

	// 2. Reserve Inventory
	ctx, reserveSpan := s.tracer.Start(ctx, "checkout.inventory.reserve")
	var reservations []*inventoryDomain.Reservation
	for _, item := range req.Items {
		res, err := s.inventoryService.Reserve(ctx, req.TenantID, item.ProductID, "checkout", item.Quantity)
		if err != nil {
			reserveSpan.RecordError(err)
			reserveSpan.End()
			// Compensation: release previous reservations
			for _, r := range reservations {
				_ = s.inventoryService.Release(ctx, req.TenantID, r.ID)
			}
			s.idempotencyFail(ctx, req.TenantID, req.IdempotencyKey)
			return nil, err
		}
		reservations = append(reservations, res)
	}
	reserveSpan.End()

	// 3. Authorize Payment
	ctx, paymentSpan := s.tracer.Start(ctx, "checkout.payment.authorize")
	authReq := paymentDomain.AuthorizeRequest{
		TenantID:       req.TenantID,
		OrderID:        "",
		Amount:         total,
		Currency:       "BRL",
		IdempotencyKey: req.IdempotencyKey,
		PaymentMethod:  req.PaymentMethod,
	}
	auth, err := s.paymentService.Authorize(ctx, "stripe", authReq)
	if err != nil {
		paymentSpan.RecordError(err)
		paymentSpan.End()
		// Compensation: release inventory
		for _, r := range reservations {
			_ = s.inventoryService.Release(ctx, req.TenantID, r.ID)
		}
		s.idempotencyFail(ctx, req.TenantID, req.IdempotencyKey)
		return nil, err
	}
	paymentSpan.End()

	// 4. Create Order
	ctx, orderSpan := s.tracer.Start(ctx, "checkout.order.create")
	order, err := s.orderService.CreateOrder(ctx, req.TenantID, req.CustomerID, orderItems)
	if err != nil {
		orderSpan.RecordError(err)
		orderSpan.End()
		// Compensation: void payment, release inventory
		_, _ = s.paymentService.Refund(ctx, "stripe", paymentDomain.RefundRequest{
			TenantID:  req.TenantID,
			CaptureID: auth.ID,
			Amount:    total,
		})
		for _, r := range reservations {
			_ = s.inventoryService.Release(ctx, req.TenantID, r.ID)
		}
		s.idempotencyFail(ctx, req.TenantID, req.IdempotencyKey)
		return nil, err
	}
	orderSpan.End()

	// 5. Capture Payment
	ctx, captureSpan := s.tracer.Start(ctx, "checkout.payment.capture")
	_, _ = s.paymentService.Capture(ctx, "stripe", paymentDomain.CaptureRequest{
		TenantID:        req.TenantID,
		AuthorizationID: auth.ID,
		Amount:          total,
		IdempotencyKey:  req.IdempotencyKey + ":capture",
	})
	captureSpan.End()

	if s.idempotency != nil && req.IdempotencyKey != "" {
		_ = s.idempotency.Complete(ctx, req.TenantID, req.IdempotencyKey, 200, []byte(fmt.Sprintf(`{"order_id":"%s"}`, order.ID)))
	}

	if err := s.orderService.UpdateStatus(ctx, req.TenantID, order.ID, orderDomain.StatusConfirmed); err != nil {
		// Non-fatal: the order was created and paid; status sync can be retried.
		_ = err
	}

	return &CheckoutResult{
		OrderID: order.ID,
		Status:  string(orderDomain.StatusConfirmed),
	}, nil
}
