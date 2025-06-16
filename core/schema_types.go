package elemental

import (
	"regexp"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type SchemaOptions struct {
	Collection              string                          // Custom collection name, if not set, the lowercase pluralized name of the model will be used.
	CollectionOptions       options.CreateCollectionOptions // Plain mongo driver collection options if you want to set them
	Database                string                          // Custom database name, if not set, the default database will be used
	Connection              string                          // Custom connection alias, if not set, the default connection will be used
	Auditing                bool                            // Whether to enable auditing for this model
	BypassSchemaEnforcement bool                            // Whether to bypass schema enforcement when creating a new document
}

type Field struct {
	Type       FieldType             // Type of the field. Can be of reflect.Kind, reflect.Type, an Elemental alias such as elemental.String or a custom reflection
	Schema     *Schema               // Defines a subschema for the field if it is a subdocument
	Required   bool                  // Whether the field is required or not when creating a new document
	Default    any                   // Default value for the field when creating a new document
	Min        float64               // Minimum value for the field when it is a number
	Max        float64               // Maximum value for the field when it is a number
	Length     int64                 // Maximum length for the field when it is a string
	Regex      *regexp.Regexp        // A regex pattern that the field must match when it is a string
	Index      *options.IndexOptions // Raw driver index options for the field. Can be used to create unique indexes, sparse indexes, etc.
	IndexOrder int                   // Sort order for the index. 1 for ascending, -1 for descending
	Ref        string                // Reference to another model if the field is a reference
	Collection string                // Collection name if the field is a reference
	IsRefID    bool                  // In development for cluster mode, don't use it yet
}
