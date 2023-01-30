package kodi

// GetSourcesRequest Files.GetSources Parameters
type GetSourcesRequest struct {
	Media string `json:"media"`
}

// GetSourcesResponse Files.GetSources Returns
type GetSourcesResponse struct {
	Limits  LimitsResult  `json:"limits"`
	Sources []*FileSource `json:"sources"`
}

type FileSource struct {
	File  string `json:"file"`
	Label string `json:"label"`
}
