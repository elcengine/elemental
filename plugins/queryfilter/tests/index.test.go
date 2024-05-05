package queryfilter

import (
	"reflect"
	"testing"
	"time"
)

var basicFilterReq = map[string]interface{}{
	"query": map[string]interface{}{
		"filter": map[string]interface{}{
			"name":       "eq(John)",
			"lastName":   "ne(Doe)",
			"middleName": "reg(.*Nathan.*)",
			"age":        "gt(20)",
			"email":      "nin(email1,email2,email3)",
			"address":    "in(address1,address2,address3)",
			"weight":     "gte(50)",
			"height":     "lt(180)",
			"birthdate":  "lte(2000-01-01)",
			"isAlive":    "exists(true)",
			"isVerified": "eq(true)",
			"isDeleted":  "false",
		},
	},
}

var basicFilterResult = map[string]interface{}{
	"name": map[string]interface{}{"$eq": "John"},
	"lastName": map[string]interface{}{"$ne": "Doe"},
	"middleName": map[string]interface{}{"$regex": "/.*Nathan.*/"},
	"age": map[string]interface{}{"$gt": 20},
	"email": map[string]interface{}{"$nin": []string{"email1", "email2", "email3"}},
	"address": map[string]interface{}{"$in": []string{"address1", "address2", "address3"}},
	"weight": map[string]interface{}{"$gte": 50},
	"height": map[string]interface{}{"$lt": 180},
	"birthdate": map[string]interface{}{"$lte": time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)},
	"isAlive": map[string]interface{}{"$exists": true},
	"isVerified": map[string]interface{}{"$eq": true},
	"isDeleted": false,
}

func TestFilterQuery(t *testing.T) {
	   basicFilterReq := make(map[string]interface{})
    basicFilterResult := make(map[string]interface{})

    MongooseFilterQuery(basicFilterReq, basicFilterResult)

    if !reflect.DeepEqual(basicFilterReq["query"], basicFilterResult) {
        t.Errorf("Filtering failed. Expected: %v, Got: %v", basicFilterResult, basicFilterReq["query"])
    }
}

