package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	TimeReflectType    = reflect.TypeOf(time.Time{})
	TimePtrReflectType = reflect.TypeOf(&time.Time{})
	ByteReflectType    = reflect.TypeOf(uint8(0))
)

func ExtractTable[T any](rows []*T, tagName string) [][]string {
	if len(rows) == 0 {
		return nil
	}

	buf := make([][]string, len(rows)+1)

	for i := 0; i < len(rows); i++ {
		row := *rows[i]
		r := reflect.ValueOf(&row).Elem()

		if i == 0 {
			buf[0] = getColumnNames(r, tagName)
		}

		buf[i+1] = getCellValues(r, tagName)
	}

	return buf
}

func getColumnNames(r reflect.Value, tagName string) []string {
	values := make([]string, 0)
	rt := r.Type()

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		value := field.Tag.Get(tagName)
		split := strings.Split(value, ";")

		if value == "-" {
			continue
		}
		//if value == "" {
		//	value = ""
		//}
		if len(split) > 1 {
			value = split[0]
		}

		values = append(values, value)
	}

	return values
}

func getCellValues(r reflect.Value, tagName string) []string {
	values := make([]string, 0)
	rt := r.Type()

	for j := 0; j < rt.NumField(); j++ {
		field := rt.Field(j)
		tag := field.Tag.Get(tagName)
		if tag == "-" {
			continue
		}

		tags := strings.Split(tag, ";")
		if len(tags) > 1 && tags[1] != "" {
			nr, nf := getNestedStructField(r, field, tags[1], tagName)
			values = append(values, serializeValue(nr, nf.Name, tagName))
			continue
		}

		values = append(values, serializeValue(r, field.Name, tagName))
	}

	return values
}

func getNestedStructField(r reflect.Value, field reflect.StructField, nestedFieldName, tagName string) (reflect.Value, reflect.StructField) {
	value := reflect.Indirect(r).FieldByName(field.Name)
	rt := value.Type()

	if value.Kind() == reflect.Ptr || value.Kind() == reflect.UnsafePointer {
		rt = value.Elem().Type()
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)

		if f.Name == nestedFieldName {
			return value.Elem(), f
		}
	}

	return r, field
}

func serializeValue(r reflect.Value, fieldName, tagName string) string {
	value := reflect.Indirect(r).FieldByName(fieldName)

	switch value.Kind() {
	case reflect.Ptr, reflect.UnsafePointer:
		if value.IsNil() {
			return "nil"
		}
		return serializeStruct(value.Elem())
	case reflect.Struct:
		return serializeStruct(value)
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Invalid, reflect.Interface, reflect.Func:
		return fmt.Sprintf("%s", value.Interface())
	case reflect.Bool, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uint8, reflect.Uintptr:
		return fmt.Sprintf("%v", value.Interface())
	case reflect.String:
		return value.String()
	case reflect.Chan:
		return ""
	}

	return ""
}

func serializeStruct(r reflect.Value) string {
	if _, ok := r.Interface().(*time.Time); ok {
		return fmt.Sprintf("%s", r.Interface())
	} else if r.Type().ConvertibleTo(TimeReflectType) {
		return fmt.Sprintf("%s", r.Interface())
	} else if r.Type().ConvertibleTo(TimePtrReflectType) {
		return fmt.Sprintf("%s", r.Interface())
	}

	return fmt.Sprintf("%s", r.Interface())
}

func Transpose(table [][]string) [][]string {
	colSize := 0

	for i := 0; i < len(table); i++ {
		if colSize < len(table[i]) {
			colSize = len(table[i])
		}
	}

	transposed := make([][]string, colSize)

	for i := 0; i < len(transposed); i++ {
		transposed[i] = make([]string, len(table))
	}

	for i := 0; i < len(table); i++ {
		for j := 0; j < len(table[i]); j++ {
			transposed[j][i] = table[i][j]
		}
	}

	return transposed
}
