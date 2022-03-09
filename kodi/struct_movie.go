package kodi

type MovieDetails struct {
	MovieId       int    `json:"movieid"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originaltitle"`
	Label         string `json:"label"`
	Year          int    `json:"year"`
}
