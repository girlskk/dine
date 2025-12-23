package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
	Create(ctx context.Context, order *Order) error
	Update(ctx context.Context, order *Order) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
}

type OrderInteractor interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	Get(ctx context.Context, id uuid.UUID) (*Order, error)
	Update(ctx context.Context, order *Order) (*Order, error)
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, params OrderListParams) ([]*Order, int, error)
}

type Order struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt int64     `json:"deleted_at"`

	MerchantID string `json:"merchant_id"`
	StoreID    string `json:"store_id"`

	BusinessDate string `json:"business_date"`
	ShiftNo      string `json:"shift_no"`
	OrderNo      string `json:"order_no"`

	OrderType     string          `json:"order_type"`
	OriginOrderID string          `json:"origin_order_id"`
	Refund        json.RawMessage `json:"refund"`

	OpenedAt    *time.Time `json:"opened_at"`
	PlacedAt    *time.Time `json:"placed_at"`
	PaidAt      *time.Time `json:"paid_at"`
	CompletedAt *time.Time `json:"completed_at"`

	OpenedBy string `json:"opened_by"`
	PlacedBy string `json:"placed_by"`
	PaidBy   string `json:"paid_by"`

	DiningMode        string `json:"dining_mode"`
	OrderStatus       string `json:"order_status"`
	PaymentStatus     string `json:"payment_status"`
	FulfillmentStatus string `json:"fulfillment_status"`
	TableStatus       string `json:"table_status"`

	TableID       string `json:"table_id"`
	TableName     string `json:"table_name"`
	TableCapacity int    `json:"table_capacity"`
	GuestCount    int    `json:"guest_count"`

	MergedToOrderID string     `json:"merged_to_order_id"`
	MergedAt        *time.Time `json:"merged_at"`

	Store   json.RawMessage `json:"store"`
	Channel json.RawMessage `json:"channel"`
	Pos     json.RawMessage `json:"pos"`
	Cashier json.RawMessage `json:"cashier"`

	Member   json.RawMessage `json:"member"`
	Takeaway json.RawMessage `json:"takeaway"`

	Cart            json.RawMessage `json:"cart"`
	Products        json.RawMessage `json:"products"`
	Promotions      json.RawMessage `json:"promotions"`
	Coupons         json.RawMessage `json:"coupons"`
	TaxRates        json.RawMessage `json:"tax_rates"`
	Fees            json.RawMessage `json:"fees"`
	Payments        json.RawMessage `json:"payments"`
	RefundsProducts json.RawMessage `json:"refunds_products"`
	Amount          json.RawMessage `json:"amount"`
}

type OrderListParams struct {
	MerchantID string
	StoreID    string

	BusinessDate string
	OrderNo      string
	OrderType    string

	OrderStatus   string
	PaymentStatus string

	Page int
	Size int
}
