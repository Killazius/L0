package validate

import (
	"errors"
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

// TODO: дописать возможные проверки данных
func Order(order *domain.Order) error {
	valid := validator.New()
	if err := registerCustomValidations(valid); err != nil {
		return err
	}
	err := valid.Struct(order)
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			return err
		}
	}
	if order.OrderUID == "" {
		return errors.New("order uid required")
	}
	return nil

}
