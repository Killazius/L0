package validate

import (
	"github.com/Killazius/L0/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func registerCustomValidations(validate *validator.Validate) error {
	err := validate.RegisterValidation("decimal", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if decimalValue, ok := field.Interface().(decimal.Decimal); ok {
			return decimalValue.GreaterThanOrEqual(decimal.Zero)
		}

		return false
	})

	if err != nil {
		return err
	}

	return nil
}

func Order(order *domain.Order) error {
	valid := validator.New()
	if err := registerCustomValidations(valid); err != nil {
		return err
	}
	return valid.Struct(order)
}
