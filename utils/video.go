package utils

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fmt"
	"github.com/fengqi/lrace"
	"github.com/spf13/cast"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	// todo 可配置
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
		"mov",
		"rmvb",
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
		".!qB",
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
	delimiterExecute = []string{ // todo 使用单独维护的音频编码、视频编码、制作组等
		"WEB-DL",
		"DDP5.1",
		"DDP 5.1",
		"DDP.5.1",
		"H.265",
		"H265",
		"BLU-RAY",
		"MA5.1",
		"MA 5.1",
		"MA.5.1",
		"MA7.1",
		"MA 7.1",
		"MA.7.1",
		"DTS-HD",
		"HDR",
		"SDR",
		"DV",
	}
	channel = []string{
		"OAD",
		"OVA",
		"BD",
		"DVD",
		"SP",
	}
	videoCoding = []string{
		"H.265",
		"H265",
		"H.264",
		"H264",
		"H.263",
		"H.261",
		"x265",
		"x264",
		"AVC",
		"MPEG",
		"av1",
		"HEVC",
	}
	audioCoding = []string{
		"ac3",
		"aac",
		"dts",
		"dts-hd",
		"e-ac-3",
		"ddp 5.1",
		"ddp5.1",
	}
	dynamicRange = []string{
		"hdr",
		"sdr",
		"dv",
	}
	crew = []string{
		"ADWeb",
		"Audies",
		"ADE",
		"ADAudio",
		"CMCT",
		"CMCTA",
		"CMCTV",
		"Oldboys",
		"GTR",
		"OurBits",
		"OurTV",
		"iLoveTV",
		"iLoveHD",
		"MTeam",
		"MWeb",
		"BMDru",
		"QHstudio",
		"HDCTV",
		"HDArea",
		"HDAccess",
		"WiKi",
		"TTG",
		"CHD",
		"beAst",
		"DTime",
		"HHWEB",
		"NoVA",
		"NoPA",
		"NoXA",
		"HDSky",
		"HDS",
		"HDSTV",
		"HDSWEB",
		"HDSPad",
		"HDS3D",
		"HDHome",
		"HDH",
		"HDHTV",
		"HDHWEB",
		"AGSVPT",
		"AGSVWEB",
	}
	videoMap        = map[string]struct{}{}
	sourceMap       = map[string]struct{}{}
	studioMap       = map[string]struct{}{}
	delimiterMap    = map[string]struct{}{}
	channelMap      = map[string]struct{}{}
	videoCodingMap  = map[string]struct{}{}
	audioCodingMap  = map[string]struct{}{}
	dynamicRangeMap = map[string]struct{}{}
	crewMap         = map[string]struct{}{}

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
		"十": 10,
	}
	chsNumberUnit = map[string]int{
		"十": 10,
		"百": 100,
		"千": 1000,
		"万": 10000,
		"亿": 100000000,
	}
	chsMatch        *regexp.Regexp
	chsSeasonMatch  *regexp.Regexp
	chsEpisodeMatch *regexp.Regexp

	episodeMatch       *regexp.Regexp
	episodeMatchAlone  *regexp.Regexp
	collectionMatch    *regexp.Regexp
	subEpisodesMatch   *regexp.Regexp
	yearRangeLikeMatch *regexp.Regexp
	yearRangeMatch     *regexp.Regexp
	yearMatch          *regexp.Regexp
	formatMatch        *regexp.Regexp
	seasonMatch        *regexp.Regexp
	optionsMatch       *regexp.Regexp
	resolutionMatch    *regexp.Regexp
	seasonRangeMatch   *regexp.Regexp
	partMatch          *regexp.Regexp
	numberMatch        *regexp.Regexp
)

func init() {
	for _, item := range video {
		videoMap[item] = struct{}{}
	}

	for _, item := range source {
		sourceMap[item] = struct{}{}
	}

	for _, item := range studio {
		studioMap[strings.ToUpper(item)] = struct{}{}
	}

	for _, item := range delimiter {
		delimiterMap[strings.ToUpper(item)] = struct{}{}
	}

	for _, item := range channel {
		channelMap[strings.ToUpper(item)] = struct{}{}
	}

	for _, item := range videoCoding {
		videoCodingMap[strings.ToUpper(item)] = struct{}{}
	}

	for _, item := range audioCoding {
		audioCodingMap[strings.ToUpper(item)] = struct{}{}
	}

	for _, item := range dynamicRange {
		dynamicRangeMap[strings.ToUpper(item)] = struct{}{}
	}

	for _, item := range crew {
		crewMap[strings.ToUpper(item)] = struct{}{}
	}

	episodeMatch, _ = regexp.Compile(`(?i)((第|s|season)\s*(\d+).*?季?)?(第|e|p|ep|episode)\s*(\d+).+$`)
	episodeMatchAlone, _ = regexp.Compile(`(?i)(第|e|p|ep|episode)\s*(\d+).+$`)
	collectionMatch, _ = regexp.Compile("[sS](0|)[0-9]+-[sS](0|)[0-9]+")
	subEpisodesMatch, _ = regexp.Compile("[eE](0|)[0-9]+-[eE](0|)[0-9]+")
	yearRangeLikeMatch, _ = regexp.Compile("[12][0-9]{3}-[12][0-9]{3}")
	yearRangeMatch, _ = regexp.Compile("[12][0-9]{3}-[12][0-9]{3}")
	yearMatch, _ = regexp.Compile("^[12][0-9]{3}$")
	formatMatch, _ = regexp.Compile("([0-9]+[pPiI]|[24][kK])")
	seasonMatch, _ = regexp.Compile(`(?i)(第|s|S|Season)\s*(\d+)(季|)(.+)?$`)
	optionsMatch, _ = regexp.Compile(`\[.*?\](\.)?`)
	chsMatch, _ = regexp.Compile("(?:第|)([零一二三四五六七八九十百千万亿]+)[季|集]")
	chsSeasonMatch, _ = regexp.Compile(`(.*?)(\.|)第([0-9]+)([-至到])?([0-9]+)?季`)
	chsEpisodeMatch, _ = regexp.Compile("(?:第|)([0-9]+)集")
	resolutionMatch, _ = regexp.Compile("[0-9]{3,4}Xx*[0-9]{3,4}")
	seasonRangeMatch, _ = regexp.Compile("[sS](0|)[0-9]+-[sS](0|)[0-9]+")
	partMatch, _ = regexp.Compile("(:?.|-|_| |@)[pP]art([0-9])(:?.|-|_| |@)")
	numberMatch, _ = regexp.Compile("([0-9]+).+$")
}

// IsCollection 是否是合集，如S01-S03季
func IsCollection(name string) bool {
	return collectionMatch.MatchString(name) || yearRangeMatch.MatchString(name)
}

// IsSubEpisodes 是否是分段集，如：World.Heritage.In.China.E01-E38.2008.CCTVHD.x264.AC3.720p-CMCT
// 常见于持续更新中的
func IsSubEpisodes(name string) string {
	return subEpisodesMatch.FindString(name)
}

// IsVideo 是否是视频文件，根据后缀枚举
func IsVideo(name string) string {
	split := strings.Split(name, ".")
	if len(split) == 0 {
		return ""
	}

	suffix := strings.ToLower(split[len(split)-1])
	if _, ok := videoMap[suffix]; ok {
		return suffix
	}

	return ""
}

// IsYearRangeLike 判断并返回年范围，用于合集
func IsYearRangeLike(name string) string {
	return yearRangeLikeMatch.FindString(name)
}

// IsYearRange 判断并返回年范围，用于合集
func IsYearRange(name string) string {
	return yearRangeMatch.FindString(name)
}

// IsYear 判断是否是年份
func IsYear(name string) int {
	if !yearMatch.MatchString(name) {
		return 0
	}

	year, _ := strconv.Atoi(name)

	return year
}

// IsSeasonRange 判断并返回合集
func IsSeasonRange(name string) string {
	return seasonRangeMatch.FindString(name)
}

// IsSeason 判断并返回季，可能和名字写在一起，所以使用子串，如：黄石S01.Yellowstone.2018.1080p
func IsSeason(name string) (string, string) {
	find := seasonMatch.FindStringSubmatch(name)
	if len(find) > 0 {
		return find[0], find[2]
	}

	return name, ""
}

// IsEpisode 判断并返回集，如果文件名带数字
func IsEpisode(name string) (string, string) {
	find := episodeMatchAlone.FindStringSubmatch(name)
	if len(find) > 0 {
		return find[0], find[2]
	}
	return name, ""
}

// IsFormat 判断并返回格式，可能放在结尾，所以使用子串，如：World.Heritage.In.China.E01-E38.2008.CCTVHD.x264.AC3.720p-CMCT
func IsFormat(name string) string {
	return formatMatch.FindString(name)
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

// IsChannel 发行渠道
func IsChannel(name string) string {
	if _, ok := channelMap[strings.ToUpper(name)]; ok {
		return name
	}
	return ""
}

// IsVideoCoding 视频编码器
func IsVideoCoding(name string) string {
	if _, ok := videoCodingMap[strings.ToUpper(name)]; ok {
		return name
	}
	return ""
}

// IsAudioCoding 音频编码
func IsAudioCoding(name string) string {
	if _, ok := audioCodingMap[strings.ToUpper(name)]; ok {
		return name
	}
	return ""
}

// IsDynamicRange 动态范围
func IsDynamicRange(name string) string {
	if _, ok := dynamicRangeMap[strings.ToUpper(name)]; ok {
		return name
	}
	return ""
}

// IsCrew 制作组
func IsCrew(name string) string {
	if _, ok := crewMap[strings.ToUpper(name)]; ok {
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

func SplitTitleAlias(name string) (string, string) {
	split := strings.Split(name, " AKA ")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return name, ""
}

// MatchEpisode 匹配季和集
func MatchEpisode(name string) (int, int) {
	find := episodeMatch.FindStringSubmatch(name)
	if len(find) == 6 {
		return cast.ToInt(find[3]), cast.ToInt(find[5])
	}

	return 0, 0
}

// FilterTmpSuffix 过滤临时文件后缀，部分软件会在未完成的文件后面增加后缀
func FilterTmpSuffix(name string) string {
	if !config.Collector.FilterTmpSuffix || len(config.Collector.TmpSuffix) == 0 {
		return name
	}

	for _, tmp := range tmpSuffix {
		for _, suffix := range video {
			name = strings.Replace(name, suffix+tmp, suffix, 1)
		}
	}

	return name
}

// FilterOptionals 过滤掉可选的字符: 被中括号[]包围的
// 若是过滤完后为空，可能直接使用[]分段，尝试只过滤第一个
func FilterOptionals(name string) string {
	clear := optionsMatch.ReplaceAllString(name, "")
	if clear != "" {
		return clear
	}

	find := optionsMatch.FindStringSubmatch(name)
	if len(find) == 2 {
		clear = strings.Replace(name, find[0], "", 1)
	}

	return clear
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
	for {
		find := chsMatch.FindStringSubmatch(name)
		if len(find) == 0 {
			break
		}

		number := strconv.Itoa(CoverChsNumber(find[1]))
		name = strings.Replace(name, find[1], number, 1)
	}

	return name
}

// SeasonCorrecting 中文季纠正
func SeasonCorrecting(name string) string {
	name = ReplaceChsNumber(name)
	right := ""
	find := chsSeasonMatch.FindStringSubmatch(name)
	if len(find) == 6 {
		if find[4] == "" && find[5] == "" {
			num, err := strconv.Atoi(find[3])
			if err == nil && num > 0 {
				right = fmt.Sprintf("S%.2d", num)
			}
		} else {
			num1, err := strconv.Atoi(find[3])
			num2, err := strconv.Atoi(find[5])
			if err == nil && num1 > 0 && num2 > 0 {
				right = fmt.Sprintf("S%.2d-S%.2d", num1, num2)
			}
		}

		if right != "" {
			name = strings.Replace(name, find[0], find[1]+"."+right, 1)
		}
	}

	return name
}

// EpisodeCorrecting 中文集纠正
func EpisodeCorrecting(name string) string {
	name = ReplaceChsNumber(name)
	find := chsEpisodeMatch.FindStringSubmatch(name)
	if len(find) == 2 {
		number, err := strconv.Atoi(find[1])
		if err == nil {
			name = strings.Replace(name, find[0], fmt.Sprintf("E%02d", number), 1)
		}
	}

	return name
}

// IsResolution 分辨率
func IsResolution(name string) string {
	return resolutionMatch.FindString(name)
}

// Split 影视目录或文件名切割
func Split(name string) []string {
	return lrace.StringSplitWith(name, delimiter, delimiterExecute)
}

// MatchPart 匹配分卷
func MatchPart(name string) int {
	find := partMatch.FindStringSubmatch(name)
	if len(find) == 4 {
		num, err := strconv.Atoi(find[2])
		if err == nil {
			return num
		}
	}
	return 0
}
