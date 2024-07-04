package flat_test

import (
	"fmt"
	"os"

	flat "github.com/quartercastle/flat-reflect"
)

func Example() {
	cfg := struct {
		Host string `default:"localhost"`
		Port string `default:"8080"`
	}{}

	for _, field := range flat.Reflect(&cfg) {
		field.SetString(field.Tag.Get("default"))
	}

	fmt.Println(cfg)
	// Output: {localhost 8080}

}

func ExampleWithEnv() {
	cfg := struct {
		Host string `env:"EXAMPLE_HOST" default:"localhost"`
		Port string `env:"EXAMPLE_PORT" default:"8080"`
	}{}

	os.Setenv("EXAMPLE_HOST", "example.com")
	os.Setenv("EXAMPLE_PORT", "443")

	for _, field := range flat.Reflect(&cfg) {
		if v, ok := os.LookupEnv(field.Tag.Get("env")); ok {
			field.SetString(v)
		} else {
			field.SetString(field.Tag.Get("default"))
		}
	}

	fmt.Println(cfg)
	// Output: {example.com 443}
}
