package elemental

import (
	"github.com/creasty/defaults"
)

type Schema struct {
	Definitions map[string]Field
	Options     SchemaOptions
}

func NewSchema(definitions map[string]Field, opts SchemaOptions) Schema {
	defaults.Set(opts)
	schema :=  Schema{
		Definitions: definitions,
		Options:     opts,
	}
	schema.SyncIndexes()	
	return schema
}

// Enables timestamps with the default field names of createdAt and updatedAt.
//
// @returns void
//
// @example
//
// schema.DefaultTimestamps()
func (s Schema) DefaultTimestamps() {
	s.Timestamps(nil)
}

// Enables timestamps with custom field names.
//
// @param ts - A struct containing the custom field names.
//
// @returns void
//
// @example
//
//	schema.Timestamps(&TS{
//		CreatedAt: "created_at",
//		UpdatedAt: "updated_at",
//	})
func (s Schema) Timestamps(ts *TS) {
	defaults.Set(s.Options.Timestamps)
	s.Options.Timestamps.Enabled = true
	if ts.CreatedAt != "" {
		s.Options.Timestamps.CreatedAt = ts.CreatedAt
	}
	if ts.UpdatedAt != "" {
		s.Options.Timestamps.UpdatedAt = ts.UpdatedAt
	}
}

func (s Schema) SyncIndexes() {
	for _, definition := range s.Definitions {
		if definition.Index {
			
		}
	}
}