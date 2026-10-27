package tmdb

import (
	"slices"
	"sort"
	"strings"
)

// Image TMDB 图片信息，电影和剧集的 images 数据结构一致，共用一个类型
type Image struct {
	AspectRatio float32 `json:"aspect_ratio"`
	Height      int     `json:"height"`
	Iso6391     string  `json:"iso_639_1"`
	FilePath    string  `json:"file_path"`
	VoteAverage float32 `json:"vote_average"`
	VoteCount   int     `json:"vote_count"`
	Width       int     `json:"width"`
}

// primaryCode 取语言标签的 ISO 639-1 主码，图片接口仅支持两字母码，区域变体（如 zh-CN）对图片不生效
func primaryCode(lang string) string {
	return strings.SplitN(lang, "-", 2)[0]
}

// includeImageLanguage 构造 include_image_language 参数：
// 配置主语言、影片原始语言都带上，en 与 null（无语言标签）恒作保底，保证候选数量足够。
// language 参数对 images 是过滤器，原始语言的图不会自动包含，必须显式列出
func includeImageLanguage(language, originalLanguage string) string {
	codes := make([]string, 0, 4)
	for _, lang := range []string{language, originalLanguage, "en", "null"} {
		code := primaryCode(lang)
		if code != "" && !slices.Contains(codes, code) {
			codes = append(codes, code)
		}
	}
	return strings.Join(codes, ",")
}

// PickLogoFile 按 配置主语言 > 影片原始语言 > 票选最高 的顺序挑选一张 logo，返回其文件路径，均无效时返回空串。
// logos 会按票选降序原地排序
func PickLogoFile(logos []*Image, language, originalLanguage string) string {
	sort.SliceStable(logos, func(i, j int) bool {
		return logos[i].VoteAverage > logos[j].VoteAverage
	})

	for _, lang := range []string{language, originalLanguage} {
		code := primaryCode(lang)
		for _, logo := range logos {
			if logo.Iso6391 == code && logo.FilePath != "" {
				return logo.FilePath
			}
		}
	}

	for _, logo := range logos {
		if logo.FilePath != "" {
			return logo.FilePath
		}
	}
	return ""
}
