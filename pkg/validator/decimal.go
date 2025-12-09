package validator

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

// RegisterDecimalValidator 注册 decimal 验证器
func RegisterDecimalValidator(validate *validator.Validate) error {
	validate.RegisterCustomTypeFunc(func(field reflect.Value) any {
		if value, ok := field.Interface().(decimal.Decimal); ok {
			return value.String()
		}
		return nil
	}, decimal.Decimal{})

	if err := validate.RegisterValidation("d_nonnegative", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		return !value.IsNegative()
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("d_positive", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		return value.IsPositive()
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("d_gte", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.GreaterThanOrEqual(baseValue)
	}); err != nil {
		return err
	}

	if err := validate.RegisterValidation("d_lte", func(fl validator.FieldLevel) bool {
		data, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		value, err := decimal.NewFromString(data)
		if err != nil {
			return false
		}
		baseValue, err := decimal.NewFromString(fl.Param())
		if err != nil {
			return false
		}
		return value.LessThanOrEqual(baseValue)
	}); err != nil {
		return err
	}

	return nil
}
