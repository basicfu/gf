package mgd

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Decimal decimal.Decimal

func (d Decimal) DecodeValue(dc bson.DecodeContext, vr bson.ValueReader, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || !val.CanSet() || val.Type() != decimalType {
		return bson.ValueDecoderError{
			Name:     "decimalDecodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}
	var value decimal.Decimal
	switch vr.Type() {
	case bson.TypeDecimal128:
		dec, err := vr.ReadDecimal128()
		if err != nil {
			return err
		}
		value, err = decimal.NewFromString(dec.String())
		if err != nil {
			return err
		}
	case bson.TypeInt64: //int64也可以转为decimal
		dec, err := vr.ReadInt64()
		if err != nil {
			return err
		}
		value = decimal.NewFromInt(dec)
	default:
		return fmt.Errorf("received invalid BSON type to decode into decimal.Decimal: %s", vr.Type())
	}
	val.Set(reflect.ValueOf(value))
	return nil
}

func (d Decimal) EncodeValue(ec bson.EncodeContext, vw bson.ValueWriter, val reflect.Value) error {
	decimalType := reflect.TypeOf(decimal.Decimal{})
	if !val.IsValid() || val.Type() != decimalType {
		return bson.ValueEncoderError{
			Name:     "decimalEncodeValue",
			Types:    []reflect.Type{decimalType},
			Received: val,
		}
	}
	dec := val.Interface().(decimal.Decimal)
	dec128, err := bson.ParseDecimal128(dec.String())
	if err != nil {
		return err
	}
	return vw.WriteDecimal128(dec128)
}
