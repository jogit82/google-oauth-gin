package structs

// Movie stores retrieved movies
type Movie struct {
	Name       string  `json:"name"`
	Overview   string  `json:"overview"`
	Popularity float32 `json:"popularity"`
	VoteCount  int64   `json:"vote_count"`
}

// Show stores retrieved shows
type Show struct {
	Name       string  `json:"name"`
	Overview   string  `json:"overview"`
	Popularity float32 `json:"popularity"`
	VoteCount  int64   `json:"vote_count"`
}

// Person stores retrieved people
type Person struct {
	Name       string  `json:"name"`
	Popularity float32 `json:"popularity"`
}
