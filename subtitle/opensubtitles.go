package subtitle

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type attributes struct {
	// SubtitleID        string    `json:"subtitle_id"`
	// Language          string    `json:"language"`
	// DownloadCount     int       `json:"download_count"`
	// NewDownloadCount  int       `json:"new_download_count"`
	// HearingImpaired   bool      `json:"hearing_impaired"`
	// Hd                bool      `json:"hd"`
	// Fps               float64   `json:"fps"`
	// Votes             int       `json:"votes"`
	// Ratings           float64   `json:"ratings"`
	// FromTrusted       bool      `json:"from_trusted"`
	// ForeignPartsOnly  bool      `json:"foreign_parts_only"`
	// UploadDate        time.Time `json:"upload_date"`
	// AiTranslated      bool      `json:"ai_translated"`
	// MachineTranslated bool      `json:"machine_translated"`
	// Release           string    `json:"release"`
	// Comments          string    `json:"comments"`
	// LegacySubtitleID  int       `json:"legacy_subtitle_id"`
	// Uploader          struct {
	// 	UploaderID int    `json:"uploader_id"`
	// 	Name       string `json:"name"`
	// 	Rank       string `json:"rank"`
	// } `json:"uploader"`
	// FeatureDetails struct {
	// 	FeatureID       int    `json:"feature_id"`
	// 	FeatureType     string `json:"feature_type"`
	// 	Year            int    `json:"year"`
	// 	Title           string `json:"title"`
	// 	MovieName       string `json:"movie_name"`
	// 	ImdbID          int    `json:"imdb_id"`
	// 	TmdbID          int    `json:"tmdb_id"`
	// 	SeasonNumber    int    `json:"season_number"`
	// 	EpisodeNumber   int    `json:"episode_number"`
	// 	ParentImdbID    int    `json:"parent_imdb_id"`
	// 	ParentTitle     string `json:"parent_title"`
	// 	ParentTmdbID    int    `json:"parent_tmdb_id"`
	// 	ParentFeatureID int    `json:"parent_feature_id"`
	// } `json:"feature_details"`
	// URL          string `json:"url"`
	// RelatedLinks []struct {
	// 	Label  string `json:"label"`
	// 	URL    string `json:"url"`
	// 	ImgURL string `json:"img_url,omitempty"`
	// } `json:"related_links"`
	// MoviehashMatch bool `json:"moviehash_match"`
	Files []struct {
		FileID   int    `json:"file_id"`
		CdNumber int    `json:"cd_number"`
		FileName string `json:"file_name"`
	} `json:"files"`
}

type data struct {
	//ID         string     `json:"id"`
	//Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}
type openSubtitlesData struct {
	TotalPages int `json:"total_pages"`
	TotalCount int `json:"total_count"`
	//PerPage    int    `json:"per_page"`
	Page int    `json:"page"`
	Data []data `json:"data"`
}

type openSubtitles struct {
	apiKey       string
	apiSubtitles string
	apiDownload  string
}

func NewOpenSubtitles() *openSubtitles {
	return &openSubtitles{
		apiKey:       "HaSKP2QrF89J5xooZPU6HcZUgPfrDpFw",
		apiSubtitles: "https://api.opensubtitles.com/api/v1/subtitles?languages=ze,zh-cn,zh-tw&tmdb_id=%v",
		apiDownload:  "https://api.opensubtitles.com/api/v1/download",
	}
}

func (o *openSubtitles) GetSubtitleList(fileInfo *FileInfo) (SubtitleInfoList []*SubtitleInfo, err error) {
	url := fmt.Sprintf(o.apiSubtitles, fileInfo.TmdbId)

	b, err := o.request(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var openSubtitlesList openSubtitlesData
	if err = json.Unmarshal(b, &openSubtitlesList); err != nil {
		return nil, err
	}

	if len(openSubtitlesList.Data) == 0 {
		return nil, errors.New("no subtitle found in opensubtitles")
	}

	for _, v := range openSubtitlesList.Data {
		for _, f := range v.Attributes.Files {
			downloadInfo, err := o.getDownloadUlr(f.FileID)
			if err != nil {
				return nil, err
			}

			var sub = SubtitleInfo{
				Url:      downloadInfo.Link,
				Filename: downloadInfo.Filename,
			}
			SubtitleInfoList = append(SubtitleInfoList, &sub)
		}
	}

	return
}

type downloadInfo struct {
	Link         string    `json:"link"`
	Filename     string    `json:"file_name"`
	Requests     int       `json:"requests"`
	Remaining    int       `json:"remaining"`
	Message      string    `json:"message"`
	ResetTime    string    `json:"reset_time"`
	ResetTimeUtc time.Time `json:"reset_time_utc"`
}

func (o *openSubtitles) getDownloadUlr(fileID int) (osbDownloadInfo *downloadInfo, err error) {
	reader := strings.NewReader(fmt.Sprintf("{\n  \"file_id\": %d\n}", fileID))
	b, err := o.request(http.MethodPost, o.apiDownload, reader)
	if err != nil {
		return nil, err
	}

	osbDownloadInfo = new(downloadInfo)
	if err = json.Unmarshal(b, &osbDownloadInfo); err != nil {
		return nil, err
	}
	return
}

func (o *openSubtitles) request(method, url string, body io.Reader) ([]byte, error) {
	utils.Logger.DebugF(url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Api-Key", o.apiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")

	httpClient := &http.Client{
		// Transport: tr,
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
