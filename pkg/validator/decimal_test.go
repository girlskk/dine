package validator

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
)

func TestRegisterDecimalValidator(t *testing.T) {
	validate := validator.New()
	err := RegisterDecimalValidator(validate)
	require.NoError(t, err)

	// Test decimal validation
	t.Run("d_nonnegative", func(t *testing.T) {
		type testStruct struct {
			Value decimal.Decimal `validate:"d_nonnegative"`
		}
		dataNoPass := testStruct{
			Value: decimal.NewFromFloat(-1),
		}
		err := validate.Struct(dataNoPass)
		require.Error(t, err)

		dataPass := testStruct{
			Value: decimal.NewFromFloat(0),
		}
		err = validate.Struct(dataPass)
		require.NoError(t, err)
	})

	t.Run("d_positive", func(t *testing.T) {
		type testStruct struct {
			Value decimal.Decimal `validate:"d_positive"`
		}
		dataNoPass := testStruct{
			Value: decimal.NewFromFloat(-1),
		}
		err := validate.Struct(dataNoPass)
		require.Error(t, err)

		dataPass := testStruct{
			Value: decimal.NewFromFloat(1),
		}
		err = validate.Struct(dataPass)
		require.NoError(t, err)
	})

	t.Run("d_gte", func(t *testing.T) {
		type testStruct struct {
			Value decimal.Decimal `validate:"d_gte=10.5"`
		}
		dataNoPass := testStruct{
			Value: decimal.NewFromFloat(10.4),
		}
		err := validate.Struct(dataNoPass)
		require.Error(t, err)

		dataPass := testStruct{
			Value: decimal.NewFromFloat(10.5),
		}
		err = validate.Struct(dataPass)
		require.NoError(t, err)
	})

	t.Run("d_lte", func(t *testing.T) {
		type testStruct struct {
			Value decimal.Decimal `validate:"d_lte=10.5"`
		}
		dataNoPass := testStruct{
			Value: decimal.NewFromFloat(10.6),
		}
		err := validate.Struct(dataNoPass)
		require.Error(t, err)

		dataPass := testStruct{
			Value: decimal.NewFromFloat(10.5),
		}
		err = validate.Struct(dataPass)
		require.NoError(t, err)
	})
}
