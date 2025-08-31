package test

import (
	"github.com/Killazius/L0/internal/domain"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/shopspring/decimal"
	"strings"
)

func GenerateOrder() *domain.Order {
	return &domain.Order{
		OrderUID:    strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: domain.Delivery{
			Name:    gofakeit.Name(),
			Phone:   "+1" + gofakeit.Numerify("##########"),
			Zip:     gofakeit.Zip(),
			City:    gofakeit.City(),
			Address: gofakeit.Address().Street,
			Region:  gofakeit.RandomString([]string{"North", "South", "West", "East"}),
			Email:   gofakeit.Email(),
		},
		Payment: domain.Payment{
			Transaction:  strings.ReplaceAll(gofakeit.UUID(), "-", ""),
			RequestID:    gofakeit.Numerify("####"),
			Currency:     gofakeit.CurrencyShort(),
			Provider:     gofakeit.RandomString([]string{"wbpay", "sberpay", "alipay"}),
			Amount:       decimal.NewFromInt(int64(gofakeit.Number(0, 2000))),
			PaymentDt:    int64(gofakeit.Number(1231231233, 1637907727)),
			Bank:         gofakeit.RandomString([]string{"alpha", "sber", "tbank", "pspb"}),
			DeliveryCost: decimal.NewFromInt(int64(gofakeit.Number(0, 2000))),
			GoodsTotal:   gofakeit.Number(0, 500),
			CustomFee:    gofakeit.Number(0, 10),
		},
		Items: []domain.Item{
			{
				ChrtID:      gofakeit.Number(100000, 999999),
				TrackNumber: "WBILMTESTTRACK",
				Price:       decimal.NewFromFloat(gofakeit.Float64()),
				Rid:         strings.ReplaceAll(gofakeit.UUID(), "-", ""),
				Name:        gofakeit.ProductName(),
				Sale:        gofakeit.Number(0, 50),
				Size:        gofakeit.Numerify("#"),
				TotalPrice:  decimal.NewFromInt(int64(gofakeit.Number(0, 5000))),
				NmID:        gofakeit.Number(100000, 999999),
				Brand:       gofakeit.Company(),
				Status:      gofakeit.HTTPStatusCode(),
			},
		},
		Locale:            gofakeit.LanguageAbbreviation(),
		InternalSignature: gofakeit.RandomString([]string{"warehouse", "", "pvz"}),
		CustomerID:        strings.ReplaceAll(gofakeit.UUID(), "-", ""),
		DeliveryService:   gofakeit.RandomString([]string{"wb", "ali", "ozon"}),
		ShardKey:          gofakeit.Numerify("##"),
		SmID:              gofakeit.Number(0, 100),
		DateCreated:       gofakeit.Date(),
		OofShard:          gofakeit.Numerify("#"),
	}
}
