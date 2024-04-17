package elemental

import (
	"elemental/utils"
	"reflect"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TS struct {
	CreatedAt string
	UpdatedAt string
}

type SchemaTimestamps struct {
	Enabled   bool
	CreatedAt string `default:"createdAt"`
	UpdatedAt string `default:"updatedAt"`
}

func (ts *SchemaTimestamps) WithDefaults() {
	e_utils.SetDefaults(ts)
}

type SchemaOptions struct {
	Collection string
	Database   string
	Connection string
	Timestamps SchemaTimestamps
}

type Field struct {
	Disabled bool
	Type     reflect.Kind
	Required bool
	Default  any
	Min      float64
	Max      float64
	Length   int64
	Regex    string
	Index    options.IndexOptions
	IndexOrder int
}
