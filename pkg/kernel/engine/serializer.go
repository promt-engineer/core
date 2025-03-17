package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"gorm.io/gorm/schema"
)

const SpinSerializerName = "spin"
const RestoringIndexesSerializerName = "restoring"

func InitSerializer(factory SpinFactory) {
	spinFactory = factory
}

func init() {
	schema.RegisterSerializer(SpinSerializerName, &SpinEngineSerializerPure{})
	schema.RegisterSerializer(RestoringIndexesSerializerName, &RestoringIndexesEngineSerializerPure{})
}

var (
	spinFactory SpinFactory
)

type SpinEngineSerializerPure struct {
}
type RestoringIndexesEngineSerializerPure struct {
}

func (SpinEngineSerializerPure) Scan(ctx context.Context,
	field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	var (
		spin Spin
	)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
		}

		spin, err = spinFactory.UnmarshalJSONSpin(bytes)
	}

	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(spin))

	return err
}

func (SpinEngineSerializerPure) Value(ctx context.Context,
	field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}

func (RestoringIndexesEngineSerializerPure) Scan(ctx context.Context,
	field *schema.Field, dst reflect.Value, dbValue interface{}) (err error) {
	var (
		restoringIndexes RestoringIndexes
	)

	if dbValue != nil {
		var bytes []byte
		switch v := dbValue.(type) {
		case []byte:
			bytes = v
		case string:
			bytes = []byte(v)
		default:
			return fmt.Errorf("failed to unmarshal JSONB value: %#v", dbValue)
		}

		restoringIndexes, err = spinFactory.UnmarshalJSONRestoringIndexes(bytes)
	}

	field.ReflectValueOf(ctx, dst).Set(reflect.ValueOf(restoringIndexes))

	return err
}

func (RestoringIndexesEngineSerializerPure) Value(ctx context.Context,
	field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error) {
	return json.Marshal(fieldValue)
}
