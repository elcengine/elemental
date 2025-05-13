package elemental

type PaginateResult[T any] struct {
	Docs       []T    `json:"docs"`       // The documents returned by the query
	TotalDocs  int64  `json:"totalDocs"`  // The total number of documents which match the query before pagination
	Page       int64  `json:"page"`       // The current page number
	Limit      int64  `json:"limit"`      // The number of documents per page
	TotalPages int64  `json:"totalPages"` // The total number of pages which match the query
	NextPage   *int64 `json:"nextPage"`   // The next page number if there is one
	PrevPage   *int64 `json:"prevPage"`   // The previous page number if there is one
	HasPrev    bool   `json:"hasPrev"`    // Whether there is a previous page or not
	HasNext    bool   `json:"hasNext"`    // Whether there is a next page or not
}

type facetResult[T any] struct {
	Docs  []T                `bson:"docs"`
	Count []map[string]int64 `bson:"count"`
}
