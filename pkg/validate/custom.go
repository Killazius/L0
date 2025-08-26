package validate

import (
	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

func RegisterCustomValidations(validate *validator.Validate) error {
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
