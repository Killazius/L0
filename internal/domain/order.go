package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required,alphanum"`
	TrackNumber       string    `json:"track_number" validate:"required"`
	Entry             string    `json:"entry" validate:"required,alpha"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" validate:"required,alpha,len=2"`
	InternalSignature string    `json:"internal_signature" validate:"omitempty"`
	CustomerID        string    `json:"customer_id" validate:"required,alphanum"`
	DeliveryService   string    `json:"delivery_service" validate:"required,alpha"`
	ShardKey          string    `json:"shardkey" validate:"required,alphanum"`
	SmID              int       `json:"sm_id" validate:"required,min=0"`
	DateCreated       time.Time `json:"date_created" validate:"required"`
	OofShard          string    `json:"oof_shard" validate:"required,alphanum"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required,e164"`
	Zip     string `json:"zip" validate:"required,numeric"`
	City    string `json:"city" validate:"required"`
	Address string `json:"address" validate:"required"`
	Region  string `json:"region" validate:"required,alphanum"`
	Email   string `json:"email" validate:"required,email"`
}

type Payment struct {
	Transaction  string          `json:"transaction" validate:"required,alphanum"`
	RequestID    string          `json:"request_id" validate:"omitempty,alphanum"`
	Currency     string          `json:"currency" validate:"required,alpha,len=3"`
	Provider     string          `json:"provider" validate:"required,alpha"`
	Amount       decimal.Decimal `json:"amount" validate:"required,decimal"`
	PaymentDt    int64           `json:"payment_dt" validate:"required,min=0"`
	Bank         string          `json:"bank" validate:"required,alpha"`
	DeliveryCost decimal.Decimal `json:"delivery_cost" validate:"required,decimal"`
	GoodsTotal   int             `json:"goods_total" validate:"required,min=0"`
	CustomFee    int             `json:"custom_fee" validate:"min=0"`
}

type Item struct {
	ChrtID      int             `json:"chrt_id" validate:"required,min=0"`
	TrackNumber string          `json:"track_number" validate:"required"`
	Price       decimal.Decimal `json:"price" validate:"required,decimal"`
	Rid         string          `json:"rid" validate:"required,alphanum"`
	Name        string          `json:"name" validate:"required"`
	Sale        int             `json:"sale" validate:"min=0,max=100"`
	Size        string          `json:"size" validate:"required"`
	TotalPrice  decimal.Decimal `json:"total_price" validate:"required,decimal"`
	NmID        int             `json:"nm_id" validate:"required,min=0"`
	Brand       string          `json:"brand" validate:"required"`
	Status      int             `json:"status" validate:"required,min=0"`
}
