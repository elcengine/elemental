package elemental

import "reflect"

type TS struct {
	CreatedAt string
	UpdatedAt string
}

type SchemaTimestamps struct {
	enabled bool
	createdAt string `default:"createdAt"`
	updatedAt string `default:"updatedAt"`
}

type SchemaOptions struct {
	timestamps SchemaTimestamps
}

type Field struct {
	Type     reflect.Kind
	Required bool
	Default  any
	Min      float64
	Max      float64
	Length   int64
	Regex    string
	Unique   bool
	Index    bool
	Validate string
}
