package flat_test

import (
	"fmt"

	flat "github.com/quartercastle/flat-reflect"
)

func ExampleToken_SetString() {
	cfg := struct {
		String  string             `default:"hello world"`
		Int     int                `default:"1"`
		Int8    int8               `default:"2"`
		Int16   int16              `default:"3"`
		Int32   int32              `default:"4"`
		Int64   int64              `default:"5"`
		UInt    uint               `default:"6"`
		UInt8   uint8              `default:"7"`
		UInt16  uint16             `default:"8"`
		UInt32  uint32             `default:"9"`
		UInt64  uint64             `default:"10"`
		Float32 float32            `default:"1.1"`
		Float64 float64            `default:"2.2"`
		Func    func()             `default:""`
		Slice   []string           `default:"1,2,3"`
		Array   [3]int             `default:"1,2"`
		Map     map[string]float64 `default:"a=1;b=2"`
		Struct  struct {
			Foo string `default:"bar"`
		}
	}{}

	for _, field := range flat.Reflect(&cfg) {
		field.SetString(field.Tag.Get("default"))
	}

	fmt.Println(cfg)
	// Output: {hello world 1 2 3 4 5 6 7 8 9 10 1.1 2.2 <nil> [1 2 3] [1 2 0] map[a:1 b:2] {bar}}
}
