package utils

import (
	"fmt"
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
		":",
		"：",
	}
	videoMap     = map[string]struct{}{}
	sourceMap    = map[string]struct{}{}
	studioMap    = map[string]struct{}{}
	delimiterMap = map[string]struct{}{}

	chsNumber = map[string]int{
		"零": 0,
		"一": 1,
		"二": 2,
		"三": 3,
		"四": 4,
		"五": 5,
		"六": 6,
		"七": 7,
		"八": 8,
		"九": 9,
	}
	chsNumberUnit = map[string]int{
		"十": 10,
		"百": 100,
		"千": 1000,
		"万": 10000,
		"亿": 100000000,
	}
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

// SplitChsEngTitle 分离中英文名字, 不兼容中英文混编,如: 我love你
func SplitChsEngTitle(name string) (string, string) {
	name = strings.Replace(name, "[", "", -1)
	name = strings.Replace(name, "]", "", -1)
	name = strings.Replace(name, "{", "", -1)
	name = strings.Replace(name, "}", "", -1)
	name = strings.Trim(name, " ")

	//chsFind := false
	chsName := ""
	split := strings.Split(name, " ")
	for _, item := range split {
		r := []rune(item)
		//if item == "" || unicode.Is(unicode.Han, r[0]) || (chsFind && unicode.IsDigit(r[0])) {
		if item == "" || unicode.Is(unicode.Han, r[0]) {
			//chsFind = true
			chsName += item + " "
			continue
		} else {
			break
		}
	}

	chsName = strings.TrimSpace(chsName)
	engName := strings.TrimSpace(strings.Replace(name, chsName, "", 1))

	return chsName, engName
}

// MatchEpisode 匹配季和集
func MatchEpisode(name string) (string, int, int) {
	compile, err := regexp.Compile("([sS]([0-9]+))?[ ._x-]?([eEpP]([0-9]+))")
	if err != nil {
		panic(err)
	}

	find := compile.FindStringSubmatch(name)
	if len(find) != 5 {
		return "", 0, 0
	}

	season := 1
	episode := 0
	if len(find[2]) > 0 {
		s, err := strconv.Atoi(find[2])
		if err == nil {
			season = s
		}
	}

	if len(find[4]) > 0 {
		s, err := strconv.Atoi(find[4])
		if err == nil {
			episode = s
		}
	}

	se := fmt.Sprintf("s%02de%02d", season, episode)

	return se, season, episode
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
	compile, err := regexp.Compile("\\[.*?\\](\\.)?")
	if err != nil {
		Logger.ErrorF("regexp compile err: %v", err)
		return name
	}

	return compile.ReplaceAllString(name, "")
}

// CoverChsNumber 中文数字替换为阿拉伯数字
func CoverChsNumber(number string) int {
	sum := 0
	temp := 0
	runes := []rune(number)
	for i := 0; i < len(runes); i++ {
		char := string(runes[i])
		if char == "零" {
			continue
		}

		if char == "亿" || char == "万" { // 特殊的权位数字，不会再累加了，其他的十、百、千可能会继续累加，比如八百一十二万
			sum += temp * chsNumberUnit[char]
			temp = 0
		} else {
			if i+1 < len(runes) { // 还没有到最后
				nextChar := string(runes[i+1])
				if unit, ok := chsNumberUnit[nextChar]; ok { // 下一位是权位数字
					if nextChar != "亿" && nextChar != "万" {
						temp += chsNumber[char] * unit
						i++
						continue
					}
				} else { // 还没有到最后，但是下一位却不是权位数字，那自己就是权位数字，比如十二
					temp += 10
					continue
				}
			}

			temp += chsNumber[char]
		}
	}

	return sum + temp
}

// ReplaceChsNumber 替换字符里面的中文数字为阿拉伯数字
func ReplaceChsNumber(name string) string {
	compile, err := regexp.Compile("([零一二三四五六七八九十百千万亿]+)")
	if err != nil {
		Logger.ErrorF("regexp compile err: %v", err)
		return name
	}

	find := compile.FindStringSubmatch(name)
	if len(find) == 2 {
		number := strconv.Itoa(CoverChsNumber(find[1]))
		name = strings.Replace(name, find[1], number, 1)
	}

	return name
}

// FilterCorrecting 特殊字符纠正为可是识别的字符，或者过滤掉
func FilterCorrecting(name string) string {
	name = ReplaceChsNumber(name)

	compile, err := regexp.Compile("第([0-9]+)([-至到])?([0-9]+)?季")
	if err != nil {
		Logger.ErrorF("regexp compile err: %v", err)
		return name
	}

	right := ""
	find := compile.FindStringSubmatch(name)
	if len(find) == 4 {
		if find[2] == "" && find[3] == "" {
			num, err := strconv.Atoi(find[1])
			if err == nil && num > 0 {
				right = fmt.Sprintf("S%.2d", num)
			}
		} else {
			num1, err := strconv.Atoi(find[1])
			num2, err := strconv.Atoi(find[3])
			if err == nil && num1 > 0 && num2 > 0 {
				right = fmt.Sprintf("S%.2d-S%.2d", num1, num2)
			}
		}
	}

	if right != "" {
		return compile.ReplaceAllString(name, right)
	}

	return name
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
