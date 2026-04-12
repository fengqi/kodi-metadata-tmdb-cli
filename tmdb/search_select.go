package tmdb

import (
	"fengqi/kodi-metadata-tmdb-cli/common/ai"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"math"
	"strings"
)

func selectMovieResult(chsTitle, engTitle string, year int, results []*SearchMoviesResults) *SearchMoviesResults {
	if len(results) == 0 {
		return nil
	}

	switch searchMode() {
	case config.AiSearchModeAiDecision:
		// AI 选不出来时回退到算法选择
		if selected := selectMovieByAI(chsTitle, engTitle, year, results); selected != nil {
			return selected
		}
		fallthrough
	case config.AiSearchModeAlgorithm:
		return selectMovieByAlgorithm(chsTitle, engTitle, year, results)
	default:
		return results[0]
	}
}

func selectShowResult(chsTitle, engTitle string, year int, results []*SearchResults) *SearchResults {
	if len(results) == 0 {
		return nil
	}

	switch searchMode() {
	case config.AiSearchModeAiDecision:
		// AI 选不出来时回退到算法选择
		if selected := selectShowByAI(chsTitle, engTitle, year, results); selected != nil {
			return selected
		}
		fallthrough
	case config.AiSearchModeAlgorithm:
		return selectShowByAlgorithm(chsTitle, engTitle, year, results)
	default:
		return results[0]
	}
}

func selectMovieByAlgorithm(chsTitle, engTitle string, year int, results []*SearchMoviesResults) *SearchMoviesResults {
	best := results[0]
	bestScore := movieScore(best, chsTitle, engTitle, year)
	for _, r := range results[1:] {
		score := movieScore(r, chsTitle, engTitle, year)
		if score > bestScore {
			best = r
			bestScore = score
		}
	}
	return best
}

func selectShowByAlgorithm(chsTitle, engTitle string, year int, results []*SearchResults) *SearchResults {
	best := results[0]
	bestScore := showScore(best, chsTitle, engTitle, year)
	for _, r := range results[1:] {
		score := showScore(r, chsTitle, engTitle, year)
		if score > bestScore {
			best = r
			bestScore = score
		}
	}
	return best
}

func selectMovieByAI(chsTitle, engTitle string, year int, results []*SearchMoviesResults) *SearchMoviesResults {
	if !ai.Enabled() {
		return nil
	}

	candidates := make([]*ai.Candidate, 0, len(results))
	for _, r := range results {
		candidates = append(candidates, &ai.Candidate{
			Id:            r.Id,
			Title:         r.Title,
			OriginalTitle: r.OriginalTitle,
			Year:          parseYear(r.ReleaseDate),
			Popularity:    float64(r.Popularity),
			VoteAverage:   float64(r.VoteAverage),
		})
	}

	id, err := ai.ChooseCandidate("movie", &ai.ParseResult{
		Title:    strings.TrimSpace(chsTitle + " " + engTitle),
		ChsTitle: chsTitle,
		EngTitle: engTitle,
		Year:     year,
	}, candidates)
	if err != nil {
		return nil
	}

	for _, r := range results {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func selectShowByAI(chsTitle, engTitle string, year int, results []*SearchResults) *SearchResults {
	if !ai.Enabled() {
		return nil
	}

	candidates := make([]*ai.Candidate, 0, len(results))
	for _, r := range results {
		candidates = append(candidates, &ai.Candidate{
			Id:            r.Id,
			Title:         r.Name,
			OriginalTitle: r.OriginalName,
			Year:          parseYear(r.FirstAirDate),
			Popularity:    float64(r.Popularity),
			VoteAverage:   float64(r.VoteAverage),
		})
	}

	id, err := ai.ChooseCandidate("tv", &ai.ParseResult{
		Title:    strings.TrimSpace(chsTitle + " " + engTitle),
		ChsTitle: chsTitle,
		EngTitle: engTitle,
		Year:     year,
	}, candidates)
	if err != nil {
		return nil
	}

	for _, r := range results {
		if r.Id == id {
			return r
		}
	}
	return nil
}

func movieScore(r *SearchMoviesResults, chsTitle, engTitle string, year int) float64 {
	// 标题优先 年份次之 评分和热度只做微调
	score := 0.0
	score += titleScore(r.Title, chsTitle, engTitle) * 3
	score += titleScore(r.OriginalTitle, chsTitle, engTitle) * 2
	score += yearScore(parseYear(r.ReleaseDate), year) * 2
	score += float64(r.VoteAverage) / 10
	score += math.Min(float64(r.Popularity)/100, 1)
	return score
}

func showScore(r *SearchResults, chsTitle, engTitle string, year int) float64 {
	// 标题优先 年份次之 评分和热度只做微调
	score := 0.0
	score += titleScore(r.Name, chsTitle, engTitle) * 3
	score += titleScore(r.OriginalName, chsTitle, engTitle) * 2
	score += yearScore(parseYear(r.FirstAirDate), year) * 2
	score += float64(r.VoteAverage) / 10
	score += math.Min(float64(r.Popularity)/100, 1)
	return score
}

func titleScore(target, chsTitle, engTitle string) float64 {
	// 先做归一化 再比较中英文标题
	target = normalizeTitle(target)
	if target == "" {
		return 0
	}

	chsTitle = normalizeTitle(chsTitle)
	engTitle = normalizeTitle(engTitle)
	score := 0.0

	if chsTitle != "" {
		if target == chsTitle {
			score += 1
		} else if strings.Contains(target, chsTitle) || strings.Contains(chsTitle, target) {
			score += 0.5
		}
	}
	if engTitle != "" {
		if target == engTitle {
			score += 1
		} else if strings.Contains(target, engTitle) || strings.Contains(engTitle, target) {
			score += 0.5
		}
	}
	return score
}

func yearScore(candidateYear, queryYear int) float64 {
	// 年份允许小范围偏差 防止上下年误差
	if queryYear <= 0 || candidateYear <= 0 {
		return 0
	}

	diff := math.Abs(float64(candidateYear - queryYear))
	switch {
	case diff == 0:
		return 1
	case diff == 1:
		return 0.5
	case diff <= 2:
		return 0.25
	default:
		return -0.25
	}
}

func parseYear(date string) int {
	if len(date) < 4 {
		return 0
	}
	year := 0
	for i := 0; i < 4; i++ {
		c := date[i]
		if c < '0' || c > '9' {
			return 0
		}
		year = year*10 + int(c-'0')
	}
	return year
}

func normalizeTitle(s string) string {
	// 统一大小写和分隔符 降低命名格式差异影响
	s = strings.ToLower(strings.TrimSpace(s))
	s = strings.ReplaceAll(s, ".", " ")
	s = strings.ReplaceAll(s, "_", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.Join(strings.Fields(s), " ")
	return s
}

func searchMode() int {
	// 未启用 AI 时固定回到首条结果模式
	if config.Ai == nil || !config.Ai.Enable {
		return config.AiSearchModeFirstResult
	}
	return config.Ai.SearchMode
}
