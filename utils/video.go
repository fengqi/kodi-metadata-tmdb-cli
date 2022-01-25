package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	video = []string{
		"mkv",
		"mp4",
		"ts",
		"avi",
		"wmv",
		"m4v",
		"flv",
		"webm",
		"mpeg",
		"mpg",
		"3gp",
		"3gpp",
		"ts",
		"iso",
	}
	source = []string{
		"web-dl",
		"blu-ray",
		"bluray",
		"hdtv",
		"cctvhd",
	}
	studio = []string{
		"hmax",
		"netflix",
		"funimation",
		"amzn",
		"hulu",
		"kktv",
		"crunchyroll",
		"bbc",
	}
	tmpSuffix = []string{
		".part",
		".!qb",
		".!ut",
	}
	delimiter = []string{
		"-",
		".",
		",",
		"_",
		" ",
		"[",
		"]",
		"(",
		")",
		"{",
		"}",
		"@",
	}
	videoMap     = map[string]struct{}{}
	sourceMap    = map[string]struct{}{}
	studioMap    = map[string]struct{}{}
	delimiterMap = map[string]struct{}{}
)

func init() {
	for _, item := range video {
		videoMap[item] = struct{}{}
	}

	for _, item := range source {
		sourceMap[item] = struct{}{}
	}

	for _, item := range studio {
		studioMap[item] = struct{}{}
	}

	for _, item := range delimiter {
		delimiterMap[item] = struct{}{}
	}
}

// IsCollection 是否是合集，如S01-S03季
func IsCollection(name string) bool {
	ok, err := regexp.MatchString("[sS](0|)[0-9]+-[sS](0|)[0-9]+", name)
	return ok && err == nil
}

// IsSubEpisodes 是否是分段集，如：World.Heritage.In.China.E01-E38.2008.CCTVHD.x264.AC3.720p-CMCT
// 常见于持续更新中的
func IsSubEpisodes(name string) bool {
	ok, err := regexp.MatchString("[eE](0|)[0-9]+-[eE](0|)[0-9]+", name)
	return ok && err == nil
}

// IsVideo 是否是视频文件，根据后缀枚举
func IsVideo(name string) string {
	split := strings.Split(name, ".")
	if len(split) == 0 {
		return ""
	}

	suffix := split[len(split)-1]
	if _, ok := videoMap[suffix]; ok {
		return suffix
	}

	return ""
}

// IsYearRangeLike 判断并返回年范围，用于合集
func IsYearRangeLike(name string) string {
	compile, err := regexp.Compile("[12][0-9]{3}-[12][0-9]{3}")
	if err != nil {
		return ""
	}

	return compile.FindString(name)
}

// IsYearRange 判断并返回年范围，用于合集
func IsYearRange(name string) string {
	compile, err := regexp.Compile("^[12][0-9]{3}-[12][0-9]{3}$")
	if err != nil {
		return ""
	}

	return compile.FindString(name)
}

// IsYear 判断是否是年份
func IsYear(name string) int {
	ok, err := regexp.MatchString("^[12][0-9]{3}$", name)
	if !ok || err != nil {
		return 0
	}

	year, _ := strconv.Atoi(name)

	return year
}

// IsSeasonRange 判断并返回合集
func IsSeasonRange(name string) string {
	compile, err := regexp.Compile("[sS](0|)[0-9]+-[sS](0|)[0-9]+")
	if err != nil {
		return ""
	}

	return compile.FindString(name)
}

// IsSeason 判断并返回季，可能和名字写在一起，所以使用子串，如：黄石S01.Yellowstone.2018.1080p
func IsSeason(name string) string {
	compile, err := regexp.Compile("[sS](0|)[0-9]+")
	if err != nil {
		return ""
	}

	return compile.FindString(name)
}

// IsFormat 判断并返回格式，可能放在结尾，所以使用子串，如：World.Heritage.In.China.E01-E38.2008.CCTVHD.x264.AC3.720p-CMCT
func IsFormat(name string) string {
	compile, err := regexp.Compile("([0-9]+[pPiI]|[24][kK])")
	if err != nil {
		return ""
	}

	return compile.FindString(name)
}

// IsSource 片源
func IsSource(name string) string {
	if _, ok := sourceMap[strings.ToLower(name)]; ok {
		return name
	}
	return ""
}

// IsStudio 发行公司
func IsStudio(name string) string {
	if _, ok := studioMap[strings.ToLower(name)]; ok {
		return name
	}
	return ""
}

// CleanTitle 名字清理，对于中英文混编的，只保留中文或者英文
// BBC.行星.The.Planets.2019.Bluray.1080p.x265.10bit.2Audios.MNHD-FRDS
// [机智的医生生活].A.Wise.Doctor's.Life.2020.S01.Complete.NF.WEB-DL.1080p.H264.AAC-CMCTV
// 中国通史.2013.全100集.国语中字￡CMCT梦幻
func CleanTitle(name string) string {
	name = strings.Replace(name, "[", "", -1)
	name = strings.Replace(name, "]", "", -1)
	name = strings.Replace(name, "{", "", -1)
	name = strings.Replace(name, "}", "", -1)
	name = strings.Trim(name, " ")

	newName := ""
	split := strings.Split(name, " ")
	for _, item := range split {
		r := []rune(item)
		if item == "" || unicode.Is(unicode.Han, r[0]) {
			newName = ""
			continue
		}
		newName += item + " "
	}

	if newName == "" {
		newName = name
	}

	return strings.TrimSpace(newName)
}

// MatchEpisode 匹配季和集
func MatchEpisode(name string) (string, int, int) {
	compile, err := regexp.Compile("[sS][0-9]+[ ._x-]?[eE][0-9]+")
	if err != nil {
		panic(err)
	}

	se := compile.FindString(name)
	if se == "" {
		compile, err = regexp.Compile("[eE][0-9]+")
		if err != nil {
			panic(err)
		}
		se = compile.FindString(name)
	}

	se = strings.ToLower(se)
	if len(se) > 0 {
		split := strings.Split(se, "e")
		if se[0:1] == "s" {
			s, err := strconv.Atoi(split[0][1:])
			if err != nil {
				panic(err)
			}
			e, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			return se, s, e
		} else if se[0:1] == "e" {
			s := 1
			e, err := strconv.Atoi(split[1])
			if err != nil {
				panic(err)
			}
			return se, s, e
		}
	}

	return se, 0, 0
}

// FilterTmpSuffix 过滤临时文件后缀，部分软件会在未完成的文件后面增加后缀
func FilterTmpSuffix(name string) string {
	for _, tmp := range tmpSuffix {
		for _, suffix := range video {
			name = strings.Replace(name, suffix+tmp, suffix, 1)
		}
	}
	return name
}

// FilterOptionals 过滤掉可选的字符: 被中括号[]包围的
func FilterOptionals(name string) string {
	compile, err := regexp.Compile("\\[.*?\\]")
	if err != nil {
		panic(err)
	}

	return compile.ReplaceAllString(name, "")
}

// IsResolution 分辨率
func IsResolution(name string) string {
	compile, err := regexp.Compile("[0-9]{3,4}Xx*[0-9]{3,4}")
	if err != nil {
		return ""
	}

	return compile.FindString(name)
}

// Split 影视目录或文件名切割
// TODO 对于web-dl, h.264, blu-ray这样的可以不切割
func Split(name string) []string {
	runeStr := []rune(name)
	split := make([]string, 0)
	start := 0
	match := false
	lastMatch := false
	for k, v := range runeStr {
		if _, ok := delimiterMap[string(v)]; ok {
			if match {
				lastMatch = true
				subStr := string(runeStr[start:k])
				if subStr != "" {
					split = append(split, subStr)
				}
				match = false
			}
			start = k + 1
		}
		lastMatch = false
		match = true
	}

	if !lastMatch {
		subStr := string(runeStr[start:])
		if subStr != "" {
			split = append(split, subStr)
		}
	}

	return split
}
