package ai

import (
	"bytes"
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/common/httpx"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
	"net/http"
	"strings"
)

type parseMediaPromptSchema struct {
	Title      string `json:"title"`
	AliasTitle string `json:"alias_title"`
	ChsTitle   string `json:"chs_title"`
	EngTitle   string `json:"eng_title"`
	Year       string `json:"year"`
	Season     string `json:"season"`
	Episode    string `json:"episode"`
	Confidence string `json:"confidence"`
}

type parseMediaPrompt struct {
	Task         string                 `json:"task"`
	MediaType    string                 `json:"media_type"`
	Path         string                 `json:"path"`
	Filename     string                 `json:"filename"`
	Rule         any                    `json:"rule,omitempty"`
	Schema       parseMediaPromptSchema `json:"schema"`
	Requirements []string               `json:"requirements"`
}

type rankPromptSchema struct {
	Id         string `json:"id"`
	Confidence string `json:"confidence"`
	Reason     string `json:"reason"`
}

type rankCandidatesPrompt struct {
	Task         string           `json:"task"`
	MediaType    string           `json:"media_type"`
	Parsed       *ParseResult     `json:"parsed"`
	Candidates   []*Candidate     `json:"candidates"`
	Schema       rankPromptSchema `json:"schema"`
	Requirements []string         `json:"requirements"`
}

type completionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type completionThinking struct {
	Type string `json:"type"`
}

type completionRequest struct {
	Model    string              `json:"model"`
	Messages []completionMessage `json:"messages"`
	Thinking completionThinking  `json:"thinking"`
}

type completionChoiceMessage struct {
	Content string `json:"content"`
}

type completionChoice struct {
	Message completionChoiceMessage `json:"message"`
}

type completionResponse struct {
	Choices []completionChoice `json:"choices"`
}

type chooseCandidateResponse struct {
	Id         int     `json:"id"`
	Confidence float64 `json:"confidence"`
}

func Enabled() bool {
	return config.Ai != nil && config.Ai.Enable && config.Ai.BaseURL != "" && config.Ai.ApiKey != "" && config.Ai.Model != ""
}

func ParseMedia(input *ParseInput) (*ParseResult, error) {
	if !Enabled() {
		return nil, errors.New("ai disabled")
	}
	if input == nil {
		return nil, errors.New("input nil")
	}

	prompt := parseMediaPrompt{
		Task:      "parse_media_filename",
		MediaType: input.MediaType,
		Path:      input.Path,
		Filename:  input.Filename,
		Rule:      input.Rule,
		Schema: parseMediaPromptSchema{
			Title:      "string",
			AliasTitle: "string",
			ChsTitle:   "string",
			EngTitle:   "string",
			Year:       "number",
			Season:     "number",
			Episode:    "number",
			Confidence: "number(0-1)",
		},
		Requirements: []string{
			"for movie set season and episode to 0",
			"for tv if episode is present but season is missing set season to 1",
			"return strict json, no markdown",
		},
	}

	content, err := completion("You extract media metadata from filenames. Return JSON only.", prompt)
	if err != nil {
		return nil, err
	}

	out := &ParseResult{}
	if err = json.Unmarshal([]byte(content), out); err != nil {
		return nil, err
	}
	normalizeParseResult(out)
	return out, nil
}

func normalizeParseResult(out *ParseResult) {
	if out == nil {
		return
	}
	out.Title = strings.TrimSpace(out.Title)
	out.AliasTitle = strings.TrimSpace(out.AliasTitle)
	out.ChsTitle = strings.TrimSpace(out.ChsTitle)
	out.EngTitle = strings.TrimSpace(out.EngTitle)
}

func ChooseCandidate(mediaType string, parse *ParseResult, candidates []*Candidate) (int, error) {
	if !Enabled() {
		return 0, errors.New("ai disabled")
	}
	if len(candidates) == 0 {
		return 0, errors.New("candidates empty")
	}

	prompt := rankCandidatesPrompt{
		Task:       "rank_tmdb_candidates",
		MediaType:  mediaType,
		Parsed:     parse,
		Candidates: candidates,
		Schema: rankPromptSchema{
			Id:         "number",
			Confidence: "number(0-1)",
			Reason:     "string",
		},
		Requirements: []string{
			"must choose one candidate id from candidates",
			"prefer exact title/year match",
			"return strict json, no markdown",
		},
	}

	content, err := completion("You choose the best TMDB candidate by parsed media metadata. Return JSON only.", prompt)
	if err != nil {
		return 0, err
	}

	resp := &chooseCandidateResponse{}
	if err = json.Unmarshal([]byte(content), resp); err != nil {
		return 0, err
	}
	for _, c := range candidates {
		if c.Id == resp.Id {
			return resp.Id, nil
		}
	}

	return 0, errors.New("ai id not in candidates")
}

func completion(system string, user any) (string, error) {
	payload := completionRequest{
		Model: config.Ai.Model,
		Messages: []completionMessage{
			{Role: "system", Content: system},
			{Role: "user", Content: mustJSON(user)},
		},
		Thinking: completionThinking{Type: "disabled"},
	}

	b, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPost, normalizeCompletionURL(config.Ai.BaseURL), bytes.NewReader(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+config.Ai.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := httpx.NewClient(currentProxy(), config.Ai.TimeoutSeconds)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errors.New("ai completion request failed")
	}

	out := &completionResponse{}
	if err = json.Unmarshal(body, out); err != nil {
		return "", err
	}
	if len(out.Choices) == 0 {
		return "", errors.New("ai response empty")
	}

	content := strings.TrimSpace(out.Choices[0].Message.Content)
	return extractJSONObject(content), nil
}

func currentProxy() string {
	if config.Tmdb == nil {
		return ""
	}
	return config.Tmdb.Proxy
}

func normalizeCompletionURL(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimRight(url, "/")
	if strings.HasSuffix(url, "/chat/completions") {
		return url
	}
	return url + "/chat/completions"
}

func extractJSONObject(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return content
	}
	start := strings.Index(content, "{")
	end := strings.LastIndex(content, "}")
	if start >= 0 && end > start {
		return content[start : end+1]
	}
	return content
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		utils.Logger.WarningF("marshal ai prompt err: %v", err)
		return "{}"
	}
	return string(b)
}

func ParseUsable(res *ParseResult) bool {
	if res == nil {
		return false
	}
	if res.Title == "" && res.ChsTitle == "" && res.EngTitle == "" {
		return false
	}
	return res.Confidence >= config.Ai.ConfidenceThreshold
}
