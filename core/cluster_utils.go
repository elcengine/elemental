package elemental

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Function to convert primitive.D to map[string]any
func convertDocToMap(doc primitive.D) map[string]any {
	result := make(map[string]any)
	for _, elem := range doc {
		result[elem.Key] = elem.Value
	}
	return result
}
