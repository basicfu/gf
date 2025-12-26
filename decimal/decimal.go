package decimal

import (
	"github.com/shopspring/decimal"
)

var New = func(value any) decimal.Decimal {
	if value == nil {
		return decimal.Zero
	}
	switch v := value.(type) {
	case int:
		return decimal.NewFromInt(int64(v))
	case int8:
		return decimal.NewFromInt(int64(v))
	case int16:
		return decimal.NewFromInt(int64(v))
	case int32:
		return decimal.NewFromInt(int64(v))
	case int64:
		return decimal.NewFromInt(v)
	case uint:
		return decimal.NewFromInt(int64(v))
	case uint8:
		return decimal.NewFromInt(int64(v))
	case uint16:
		return decimal.NewFromInt(int64(v))
	case uint32:
		return decimal.NewFromInt(int64(v))
	case uint64:
		return decimal.NewFromInt(int64(v))
	case float32:
		return decimal.NewFromFloat32(v)
	case float64:
		return decimal.NewFromFloat(v)
	case string:
		d, err := decimal.NewFromString(v)
		if err != nil {
			return decimal.Zero
		}
		return d
	case decimal.Decimal:
		return v
	default:
		return decimal.Zero
	}
}
