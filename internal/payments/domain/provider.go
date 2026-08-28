package domain

type Provider interface {
	Authorize(ctx context.Context, request AuthorizeRequest) (Authorization, error)
	Capture(ctx context.Context, request CaptureRequest) (Capture, error)
	Refund(ctx context.Context, request RefundRequest) (Refund, error)
}

type AuthorizeRequest struct {
	TenantID      string `json:"tenant_id"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"` // cents
	Currency      string `json:"currency"`
	IdempotencyKey string `json:"idempotency_key"`
	PaymentMethod PaymentMethod `json:"payment_method"`
}

type PaymentMethod struct {
	Type       string `json:"type"`
	Token      string `json:"token"`
	Last4      string `json:"last4,omitempty"`
	ExpiryMonth string `json:"expiry_month,omitempty"`
	ExpiryYear  string `json:"expiry_year,omitempty"`
}

type Authorization struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	Provider  string `json:"provider"`
	CreatedAt string `json:"created_at"`
}

type CaptureRequest struct {
	TenantID         string `json:"tenant_id"`
	AuthorizationID  string `json:"authorization_id"`
	Amount           int64  `json:"amount"`
	IdempotencyKey   string `json:"idempotency_key"`
}

type Capture struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Amount    int64  `json:"amount"`
	CreatedAt string `json:"created_at"`
}

type RefundRequest struct {
	TenantID       string `json:"tenant_id"`
	CaptureID      string `json:"capture_id"`
	Amount         int64  `json:"amount"`
	IdempotencyKey string `json:"idempotency_key"`
}

type Refund struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	Amount    int64  `json:"amount"`
	CreatedAt string `json:"created_at"`
}
