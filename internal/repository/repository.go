package repository

import (
	"errors"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrDeliveryNotFound = errors.New("delivery not found")
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrItemsNotFound    = errors.New("items not found")
	ErrDuplicateOrder   = errors.New("duplicate order")
)
