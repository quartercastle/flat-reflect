package flat

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
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

func (token *Token) Set(value any) error {
	t := token.Value.Type()
	v := reflect.ValueOf(value)

	if t != v.Type() {
		return errors.New("not of same type")
	}

	token.Value.Set(v)
	return nil
}

func setString(ref reflect.Value, value string) error {
	t := ref.Type()
	switch t.Kind() {
	case reflect.String:
		ref.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		value, err := strconv.ParseInt(value, 0, t.Bits())
		if err != nil {
			return err
		}
		ref.SetInt(int64(value))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value, err := strconv.ParseUint(value, 0, t.Bits())
		if err != nil {
			return err
		}
		ref.SetUint(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		ref.SetBool(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(value, t.Bits())
		if err != nil {
			return err
		}
		ref.SetFloat(value)
	case reflect.Slice:
		if t.Elem().Kind() == reflect.Struct {
			// cannot handle struct initialization in lists
			return nil
		}
		values := strings.Split(value, ",")
		slice := reflect.MakeSlice(t, len(values), len(values))
		for i, v := range values {
			if err := setString(slice.Index(i), strings.TrimSpace(v)); err != nil {
				return err
			}
		}
		ref.Set(slice)
	case reflect.Array:
		values := strings.Split(value, ",")
		for i, v := range values {
			if i >= ref.Len() {
				return errors.New("array out of bounds")
			}
			if err := setString(ref.Index(i), strings.TrimSpace(v)); err != nil {
				return err
			}
		}
	case reflect.Map:
		if ref.IsNil() {
			ref.Set(reflect.MakeMap(t))
		}
		entries := strings.Split(value, ";")
		for _, entry := range entries {
			parts := strings.Split(entry, "=")
			if len(parts) == 0 {
				continue
			}
			if len(parts) < 2 {
				parts = append(parts, make([]string, len(parts)-2)...)
			}
			key, values := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			k := reflect.New(t.Key()).Elem()
			if err := setString(k, key); err != nil {
				return err
			}
			v := reflect.New(t.Elem()).Elem()
			if err := setString(v, values); err != nil {
				return err
			}
			ref.SetMapIndex(k, v)
		}
	case reflect.Struct:
	}
	return nil
}

func (token *Token) SetString(value string) error {
	return setString(token.Value, value)
}
