package queryfilter

import (
	"log"
	"net/http"
	"strings"
	"strconv"
	"./plugins/queryfilter"
	
)


func MongooseFilterQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		defer func() {
			if r := recover(); r != nil {
				log.Println("[ FilterQuery ] - Failed to parse query", r)
			}
		}()
		
		filterMap := make(map[string]string)
		for key, value := range query {
			filterMap[key] = value[0] 
		}

		reqFilter := MapFilters(filterMap)
		_ = MapFilters(filterMap)

		sortValues := query.Get("sort")
		if sortValues != "" {
			sortMap := make(map[string]interface{})
			sortPairs := strings.Split(sortValues, ",")
			for _, pair := range sortPairs {
				keyValue := strings.Split(pair, "=")
				if len(keyValue) == 2 {
					key := keyValue[0]
					dir := keyValue[1]
					if dir == "1" || dir == "-1" {
						dirInt, _ := strconv.Atoi(dir)
						sortMap[key] = dirInt
					}
				}
			}
			reqFilter["sort"] = sortMap
		}

		r.ParseForm()
		includeValues := r.Form.Get("include")
		if includeValues != "" {
			include := strings.Split(includeValues, ",")
			reqFilter["include"] = include
		}

		selectValues := r.Form.Get("select")
		if selectValues != "" {
			selectFields := strings.Split(selectValues, ",")
			reqFilter["select"] = strings.Join(selectFields, " ")
		}

		next.ServeHTTP(w, r)
	})
}