# flat-reflect

Flat reflect takes a slice, array, map or struct and turns it into a flattened map
where the key is the nested path and the value is of type `flat.Token`.
The flattened map can be used to manipulate the values of the initial provided type.
This is useful in cases where the type is initialized from multiple sources,
it could be a configuration file, the environment and default values.

### Install
```sh
go get github.com/quartercastle/flat-reflect
```

### Usage
Below is an example of how flat reflection can be used to initialize a struct
based on its StructTags. This is partically useful for configuration
initializations.
```go
// Initialising the cfg struct to nil
var cfg struct {
  Host string `env:"HOST" default:"localhost"`
  Port uint   `env:"PORT" default:"8080"`
}

// Using flat.Reflect to flatten the structure of cfg and loop over it to
// set its values based on the defined StructTag. It will set the value to either the
// environment variable or set it to the default value in this example.
for _, field := range flat.Reflect(&cfg) {
  if v, ok := os.LookupEnv(field.Tag.Get("env")); ok {
    field.SetString(v)
  } else {
    field.SetString(field.Tag.Get("default"))
  }
}

fmt.Println(cfg)
// Output: {localhost, 8080}
```

