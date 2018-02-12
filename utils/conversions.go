package utils

import (
	"strconv"

	"github.com/shopspring/decimal"
)

func ToFloat64(v interface{}) float64 {
	if v == nil {
		return 0.0
	}

	switch v.(type) {
	case float64:
		return v.(float64)
	case string:
		vStr := v.(string)
		vF, _ := strconv.ParseFloat(vStr, 64)
		return vF
	case int:
		return float64(v.(int))
	case int32:
		return float64(v.(int32))
	case int64:
		return float64(v.(int64))
	default:
		return 0.0
	}
}

func ToDecimal(v interface{}) decimal.Decimal {
	if v == nil {
		return decimal.Zero
	}

	switch v.(type) {
	case float64:
		return decimal.NewFromFloat(v.(float64))
	case string:
		d, err := decimal.NewFromString(v.(string))
		if err != nil {
			return decimal.Zero
		}
		return d
	case int:
		return decimal.NewFromFloat(float64(v.(int)))
	case int32:
		return decimal.NewFromFloat(float64(v.(int32)))
	case int64:
		return decimal.NewFromFloat(float64(v.(int64)))

	default:
		return decimal.Zero
	}
}

func ToDecimalAndRound(v interface{}, places int32) decimal.Decimal {
	if v == nil {
		return decimal.Zero
	}

	switch v.(type) {
	case float64:
		return decimal.NewFromFloat(v.(float64)).Round(places)
	case string:
		d, err := decimal.NewFromString(v.(string))
		if err != nil {
			return decimal.Zero
		}
		return d.Round(places)
	case int:
		return decimal.NewFromFloat(float64(v.(int))).Round(places)
	case int32:
		return decimal.NewFromFloat(float64(v.(int32))).Round(places)
	case int64:
		return decimal.NewFromFloat(float64(v.(int64))).Round(places)

	default:
		return decimal.Zero
	}
}

func ToInt(v interface{}) int {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case string:
		vStr := v.(string)
		vInt, _ := strconv.Atoi(vStr)
		return vInt
	case int:
		return v.(int)
	case float64:
		vF := v.(float64)
		return int(vF)
	default:
		return 0
	}
}

func ToUint64(v interface{}) uint64 {
	if v == nil {
		return 0
	}

	switch v.(type) {
	case int:
		return uint64(v.(int))
	case float64:
		return uint64((v.(float64)))
	case string:
		uV, _ := strconv.ParseUint(v.(string), 10, 64)
		return uV
	default:
		return 0
	}
}
