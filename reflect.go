package flat

import (
	"fmt"
	"reflect"
	"strings"
)

func Reflect(value any) Map {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return flatten(v)
}

func flatten(v reflect.Value, keys ...string) Map {
	tokens := Map{}
	t := v.Type()
	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			tokens[strings.Join(append(keys, fmt.Sprint(i)), ".")] = &Token{
				Value: v.Index(i),
			}

			tokens = Merge(
				tokens,
				flatten(
					v.Index(i),
					strings.Join(append(keys, fmt.Sprint(i)), "."),
				),
			)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			tokens[strings.Join(append(keys, fmt.Sprint(k)), ".")] = &Token{
				Value: v.MapIndex(k),
			}

			tokens = Merge(
				tokens,
				flatten(
					v.MapIndex(k),
					strings.Join(append(keys, fmt.Sprint(k)), "."),
				),
			)
		}
	case reflect.Struct:
		fields := reflect.VisibleFields(t)
		for i, field := range fields {
			if !field.IsExported() {
				continue
			}

			tokens[strings.Join(append(keys, fmt.Sprint(field.Name)), ".")] = &Token{
				Value: v.Field(i),
				Tag:   field.Tag,
			}

			tokens = Merge(
				tokens,
				flatten(
					v.Field(i),
					strings.Join(append(keys, field.Name), "."),
				),
			)
		}
	}

	return tokens
}
