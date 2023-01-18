package flatReflect

import (
	"errors"
	"reflect"
	"strconv"
)

type Token struct {
	Value reflect.Value
	Tag   reflect.StructTag
}

func (token Token) Type() string {
	t := token.Value.Type()
	var ident string

	if v := t.PkgPath(); v != "" {
		ident += v + "."
	}

	return build(ident, t)
}

func build(ident string, t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		ident += "*"
		t = t.Elem()
		return build(ident, t)
	}

	if t.Kind() == reflect.Slice {
		return build(ident+"[]", t.Elem())
	}

	if t.Kind() == reflect.Map {
		return build(ident+"map["+t.Key().Name()+"]", t.Elem())
	}

	return ident + t.Name()
}

func (token Token) Set(value any) error {
	t := token.Value.Type()
	v := reflect.ValueOf(value)

	if t != v.Type() {
		return errors.New("not of same type")
	}

	token.Value.Set(v)
	return nil
}

func (token Token) SetString(value string) error {
	t := token.Value.Type()
	switch t.Kind() {
	case reflect.String:
		token.Value.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(value, 0, t.Bits())
		if err != nil {
			return err
		}
		token.Value.SetInt(int64(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(value, 0, t.Bits())
		if err != nil {
			return err
		}
		token.Value.SetUint(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		token.Value.SetBool(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(value, t.Bits())
		if err != nil {
			return err
		}
		token.Value.SetFloat(value)
	case reflect.Slice:
	case reflect.Array:
	case reflect.Map:
	case reflect.Struct:
	}
	return nil
}
