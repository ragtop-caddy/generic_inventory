package api

// queryKey - Struct to hold the query key
type queryKey struct {
	key string
}

// queryVal - Struct to hold the query value
type queryVal struct {
	value []string
}

// URLQuery - Struct to hold URL query params
type URLQuery struct {
	raw string
}

// ParseQuery - Parse the raw query string into usable values
func (q *URLQuery) ParseQuery() {

}
