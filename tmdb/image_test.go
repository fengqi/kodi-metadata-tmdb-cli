package tmdb

import "testing"

func TestIncludeImageLanguage(t *testing.T) {
	cases := []struct {
		language string
		original string
		want     string
	}{
		{"zh-CN", "fr", "zh,fr,en,null"},
		{"zh-CN", "zh", "zh,en,null"},
		{"en-US", "en", "en,null"},
		{"fr-FR", "en", "fr,en,null"},
		{"zh", "", "zh,en,null"},
		{"", "", "en,null"},
	}

	for _, c := range cases {
		if got := includeImageLanguage(c.language, c.original); got != c.want {
			t.Fatalf("includeImageLanguage(%q, %q) = %q; want %q", c.language, c.original, got, c.want)
		}
	}
}

func TestPickLogoFile(t *testing.T) {
	logos := []*Image{
		{FilePath: "en.png", Iso6391: "en", VoteAverage: 10},
		{FilePath: "null.png", Iso6391: "", VoteAverage: 8},
		{FilePath: "zh.png", Iso6391: "zh", VoteAverage: 5},
	}

	// 主语言命中
	if got := PickLogoFile(logos, "zh-CN", "fr"); got != "zh.png" {
		t.Fatalf("PickLogoFile primary hit = %q; want zh.png", got)
	}

	// 主语言未命中时取原始语言
	logos = []*Image{
		{FilePath: "en.png", Iso6391: "en", VoteAverage: 10},
		{FilePath: "fr.png", Iso6391: "fr", VoteAverage: 5},
	}
	if got := PickLogoFile(logos, "zh-CN", "fr-FR"); got != "fr.png" {
		t.Fatalf("PickLogoFile original hit = %q; want fr.png", got)
	}

	// 均未命中时取票选最高
	if got := PickLogoFile(logos, "ja-JP", "ko"); got != "en.png" {
		t.Fatalf("PickLogoFile fallback = %q; want en.png", got)
	}

	// 无有效图片返回空
	empty := []*Image{{FilePath: "", Iso6391: "zh", VoteAverage: 9}}
	if got := PickLogoFile(empty, "zh-CN", ""); got != "" {
		t.Fatalf("PickLogoFile empty = %q; want \"\"", got)
	}
}
