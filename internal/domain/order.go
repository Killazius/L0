package domain

import (
	"github.com/shopspring/decimal"
	"time"
)

// Order represents an order entity
// @Description Order information
type Order struct {
	OrderUID          string    `json:"order_uid" validate:"required,alphanum" example:"b563feb7b2b84b6test"`
	TrackNumber       string    `json:"track_number" validate:"required" example:"WBILMTESTTRACK"`
	Entry             string    `json:"entry" validate:"required,alpha" example:"WBIL"`
	Delivery          Delivery  `json:"delivery" validate:"required"`
	Payment           Payment   `json:"payment" validate:"required"`
	Items             []Item    `json:"items" validate:"required,min=1,dive"`
	Locale            string    `json:"locale" validate:"required,alpha,len=2" example:"en"`
	InternalSignature string    `json:"internal_signature" validate:"omitempty" example:""`
	CustomerID        string    `json:"customer_id" validate:"required,alphanum" example:"test"`
	DeliveryService   string    `json:"delivery_service" validate:"required,alpha" example:"meest"`
	ShardKey          string    `json:"shardkey" validate:"required,alphanum" example:"9"`
	SmID              int       `json:"sm_id" validate:"required,min=0" example:"99"`
	DateCreated       time.Time `json:"date_created" validate:"required" example:"2021-11-26T06:22:19Z"`
	OofShard          string    `json:"oof_shard" validate:"required,alphanum" example:"1"`
}

type Delivery struct {
	Name    string `json:"name" validate:"required" example:"Test Testov"`
	Phone   string `json:"phone" validate:"required,e164" example:"+9720000000"`
	Zip     string `json:"zip" validate:"required,numeric" example:"2639809"`
	City    string `json:"city" validate:"required" example:"Kiryat Mozkin"`
	Address string `json:"address" validate:"required" example:"Ploshad Mira 15"`
	Region  string `json:"region" validate:"required,alphanum" example:"Kraiot"`
	Email   string `json:"email" validate:"required,email" example:"test@gmail.com"`
}

type Payment struct {
	Transaction  string          `json:"transaction" validate:"required,alphanum" example:"b563feb7b2b84b6test"`
	RequestID    string          `json:"request_id" validate:"omitempty,alphanum" example:""`
	Currency     string          `json:"currency" validate:"required,alpha,len=3" example:"USD"`
	Provider     string          `json:"provider" validate:"required,alpha" example:"wbpay"`
	Amount       decimal.Decimal `json:"amount" validate:"required,decimal" example:"1817"`
	PaymentDt    int64           `json:"payment_dt" validate:"required,min=0" example:"1637907727"`
	Bank         string          `json:"bank" validate:"required,alpha" example:"alpha"`
	DeliveryCost decimal.Decimal `json:"delivery_cost" validate:"required,decimal" example:"1500"`
	GoodsTotal   int             `json:"goods_total" validate:"required,min=0" example:"317"`
	CustomFee    int             `json:"custom_fee" validate:"min=0" example:"0"`
}

type Item struct {
	ChrtID      int             `json:"chrt_id" validate:"required,min=0" example:"9934930"`
	TrackNumber string          `json:"track_number" validate:"required" example:"WBILMTESTTRACK"`
	Price       decimal.Decimal `json:"price" validate:"required,decimal" example:"453"`
	Rid         string          `json:"rid" validate:"required,alphanum" example:"ab4219087a764ae0btest"`
	Name        string          `json:"name" validate:"required" example:"Mascaras"`
	Sale        int             `json:"sale" validate:"min=0,max=100" example:"30"`
	Size        string          `json:"size" validate:"required" example:"0"`
	TotalPrice  decimal.Decimal `json:"total_price" validate:"required,decimal" example:"317"`
	NmID        int             `json:"nm_id" validate:"required,min=0" example:"2389212"`
	Brand       string          `json:"brand" validate:"required" example:"Vivienne Sabo"`
	Status      int             `json:"status" validate:"required,min=0" example:"202"`
}
