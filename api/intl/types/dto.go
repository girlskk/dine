package types

import (
	"database/sql"

	"github.com/shopspring/decimal"
	pbdecimal "google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SQLNullTimeToProto(t sql.NullTime) *timestamppb.Timestamp {
	if t.Valid {
		return timestamppb.New(t.Time)
	}
	return nil
}

func ProtoTimeToSQLNullTime(t *timestamppb.Timestamp) sql.NullTime {
	if t.IsValid() {
		return sql.NullTime{Time: t.AsTime(), Valid: true}
	}
	return sql.NullTime{}
}

func ProtoDecimalToDecimal(d *pbdecimal.Decimal) (decimal.Decimal, error) {
	if d == nil {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(d.GetValue())
}

func DecimalToProtoDecimal(d decimal.Decimal) *pbdecimal.Decimal {
	return &pbdecimal.Decimal{Value: d.String()}
}
