package ai

type ParseInput struct {
	MediaType string `json:"media_type"`
	Path      string `json:"path"`
	Filename  string `json:"filename"`
	Rule      any    `json:"rule,omitempty"`
}

type ParseResult struct {
	Title      string  `json:"title"`
	AliasTitle string  `json:"alias_title"`
	ChsTitle   string  `json:"chs_title"`
	EngTitle   string  `json:"eng_title"`
	Year       int     `json:"year"`
	Season     int     `json:"season"`
	Episode    int     `json:"episode"`
	Confidence float64 `json:"confidence"`
}

type Candidate struct {
	Id            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Year          int     `json:"year"`
	Popularity    float64 `json:"popularity"`
	VoteAverage   float64 `json:"vote_average"`
}
