package flat_test

import (
	"fmt"
	"reflect"
	"testing"

	flat "github.com/quartercastle/flat-reflect"
)

type nested struct {
	Value string
	Map   map[string]string
	Slice []string
}

func TestReflect(t *testing.T) {
	s := struct {
		Value      string
		NilMap     map[int]int
		Map        map[int]string
		MapMap     map[string]map[string]int
		MapStruct  map[int]nested
		MapSlice   map[uint][]int
		Slice      []string
		SliceMap   []map[string]string
		SliceSlice [][]int
		Struct     nested
		Nil        *struct{}
		Func       func()
	}{
		Map: map[int]string{
			1: "test",
		},
		MapMap: map[string]map[string]int{
			"hello": map[string]int{
				"world": 1,
			},
		},
		MapStruct: map[int]nested{
			1: {},
		},
		MapSlice: map[uint][]int{
			1: {1, 2},
		},
		Slice: []string{"test"},
		SliceMap: []map[string]string{
			{"hello": "world"},
		},
		SliceSlice: [][]int{
			{1},
			{2},
		},
		Struct: nested{
			Map: map[string]string{
				"hello": "world",
			},
			Slice: []string{"hello", "world"},
		},
		Func: func() {},
	}

	expected := map[string]*flat.Token{
		"Value":              {Value: reflect.ValueOf(s.Value)},
		"NilMap":             {Value: reflect.ValueOf(s.NilMap)},
		"Map":                {Value: reflect.ValueOf(s.Map)},
		"Map.1":              {Value: reflect.ValueOf(s.Map[1])},
		"MapMap":             {Value: reflect.ValueOf(s.MapMap)},
		"MapMap.hello":       {Value: reflect.ValueOf(s.MapMap["hello"])},
		"MapMap.hello.world": {Value: reflect.ValueOf(s.MapMap["hello"]["world"])},
		"MapStruct":          {Value: reflect.ValueOf(s.MapStruct)},
		"MapStruct.1":        {Value: reflect.ValueOf(s.MapStruct[1])},
		"MapStruct.1.Value":  {Value: reflect.ValueOf(s.MapStruct[1].Value)},
		"MapStruct.1.Map":    {Value: reflect.ValueOf(s.MapStruct[1].Map)},
		"MapStruct.1.Slice":  {Value: reflect.ValueOf(s.MapStruct[1].Slice)},
		"MapSlice":           {Value: reflect.ValueOf(s.MapSlice)},
		"MapSlice.1":         {Value: reflect.ValueOf(s.MapSlice[1])},
		"MapSlice.1.0":       {Value: reflect.ValueOf(s.MapSlice[1][0])},
		"MapSlice.1.1":       {Value: reflect.ValueOf(s.MapSlice[1][1])},
		"Slice":              {Value: reflect.ValueOf(s.Slice)},
		"Slice.0":            {Value: reflect.ValueOf(s.Slice[0])},
		"SliceMap":           {Value: reflect.ValueOf(s.SliceMap)},
		"SliceMap.0":         {Value: reflect.ValueOf(s.SliceMap[0])},
		"SliceMap.0.hello":   {Value: reflect.ValueOf(s.SliceMap[0]["hello"])},
		"SliceSlice":         {Value: reflect.ValueOf(s.SliceSlice)},
		"SliceSlice.0":       {Value: reflect.ValueOf(s.SliceSlice[0])},
		"SliceSlice.0.0":     {Value: reflect.ValueOf(s.SliceSlice[0][0])},
		"SliceSlice.1":       {Value: reflect.ValueOf(s.SliceSlice[1])},
		"SliceSlice.1.0":     {Value: reflect.ValueOf(s.SliceSlice[1][0])},
		"Struct":             {Value: reflect.ValueOf(s.Struct)},
		"Struct.Value":       {Value: reflect.ValueOf(s.Struct.Value)},
		"Struct.Map":         {Value: reflect.ValueOf(s.Struct.Map)},
		"Struct.Map.hello":   {Value: reflect.ValueOf(s.Struct.Map["hello"])},
		"Struct.Slice":       {Value: reflect.ValueOf(s.Struct.Slice)},
		"Struct.Slice.0":     {Value: reflect.ValueOf(s.Struct.Slice[0])},
		"Struct.Slice.1":     {Value: reflect.ValueOf(s.Struct.Slice[1])},
		"Nil":                {Value: reflect.ValueOf(s.Nil)},
		"Func":               {Value: reflect.ValueOf(s.Func)},
	}

	tokens := flat.Reflect(s)

	if len(expected) != len(tokens) {
		t.Errorf("expected length of %d; go %d", len(expected), len(tokens))
	}

	for key := range expected {
		t.Run(key, func(t *testing.T) {
			e, a := expected[key], tokens[key]

			if fmt.Sprint(e.Value) != fmt.Sprint(a.Value) {
				t.Errorf("expected %s; got %s", e.Value, a.Value)
			}
		})
	}
}

func BenchmarkReflect(b *testing.B) {
	s := struct {
		Value  string
		Map    map[int]string
		Slice  []string
		Struct nested
		Nil    *struct{}
	}{
		Map: map[int]string{
			1: "test",
		},
		Slice: []string{"test"},
		Struct: nested{
			Map: map[string]string{
				"hello": "world",
			},
			Slice: []string{"hello", "world"},
		},
	}

	for i := 0; i < b.N; i++ {
		_ = flat.Reflect(s)
	}
}
