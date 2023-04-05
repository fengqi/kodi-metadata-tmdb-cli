package utils

import (
	"testing"
)

func TestMatchEpisode(t *testing.T) {
	cases := map[string][]int{
		"[堕落].The.Fall.2013.S02.E03.Complete.BluRay.720p.x264.AC3-CMCT.mkv":                 {2, 3},
		"[堕落].The.Fall.2013.S10.E12.Complete.BluRay.720p.x264.AC3-CMCT.mkv":                 {10, 12},
		"[堕落].The.Fall.2013.S120.E132.Complete.BluRay.720p.x264.AC3-CMCT.mkv":               {120, 132},
		"Agent.Carter.S02E01.1080p.BluRay.DD5.1.x264-HDS.mkv":                               {2, 1},
		"[壹高清]21点灵.Leave No Soul Behind.Ep01.HDTV.1080p.H264-OneHD.ts":                      {1, 1},
		"Kimetsu.no.Yaiba.Yuukaku-hen.E01.2021.Crunchyroll.WEB-DL.1080p.x264.AAC-HDCTV.mkv": {1, 1},
		"Overlord IV - 09.mp4":  {1, 9},
		"宝贝揪揪 第3季 第10集.mp4":     {3, 10},
		"宝贝揪揪 第9集.mp4":          {1, 9},
		"Gannibal.E01.2022":     {1, 1},
		"Gannibal E01 2022":     {1, 1},
		"Gannibal S01 E01 2022": {1, 1},
	}
	for name, cse := range cases {
		_, s, e := MatchEpisode(name)
		if s != cse[0] {
			t.Errorf("MatchEpisode(%s)\n season %d; expected %d", name, s, cse[0])
		}
		if e != cse[1] {
			t.Errorf("MatchEpisode(%s)\n episode %d; expected %d", name, e, cse[1])
		}

	}
}

func TestIsFormat(t *testing.T) {
	unit := map[string]string{
		"720":        "",
		"a720p":      "720p",
		"720P":       "720P",
		"1080p":      "1080p",
		"1080P":      "1080P",
		"2k":         "2k",
		"2K":         "2K",
		"4K":         "4K",
		"720p-CMCT":  "720p",
		"-720p-CMCT": "720p",
	}

	for k, v := range unit {
		actual := IsFormat(k)
		if actual != v {
			t.Errorf("isFormat(%s) = %s; expected %s", k, actual, v)
		}
	}
}

func TestIsSeason(t *testing.T) {
	unit := map[string]string{
		"s01":  "s01",
		"S01":  "S01",
		"s1":   "s1",
		"S1":   "S1",
		"S100": "S100",
		"4K":   "",
		"Fall.in.Love.2021.WEB-DL.4k.H265.10bit.AAC-HDCTV FallinLove ": "",
		"Hawkeye.2021S01.Never.Meet.Your.Heroes.2160p":                 "S01",
	}

	for k, v := range unit {
		actual := IsSeason(k)
		if actual != v {
			t.Errorf("isSeason(%s) = %s; expected %s", k, actual, v)
		}
	}
}

func TestSplit(t *testing.T) {
	unit := map[string][]string{
		"[梦蓝字幕组]Crayonshinchan 蜡笔小新[1105][2021.11.06][AVC][1080P][GB_JP][MP4]V2.mp4": {
			"梦蓝字幕组",
			"Crayonshinchan",
			"蜡笔小新",
			"1105",
			"2021",
			"11",
			"06",
			"AVC",
			"1080P",
			"GB",
			"JP",
			"MP4",
			"V2",
			"mp4",
		},
		"The Last Son 2021.mkv": {
			"The",
			"Last",
			"Son",
			"2021",
			"mkv",
		},
		"Midway 2019 2160p CAN UHD Blu-ray HEVC DTS-HD MA 5.1-THDBST@HDSky.nfo": {
			"Midway",
			"2019",
			"2160p",
			"CAN",
			"UHD",
			"Blu-ray",
			"HEVC",
			"DTS-HD",
			"MA 5.1",
			"THDBST",
			"HDSky",
			"nfo",
		},
	}

	for k, v := range unit {
		actual := Split(k)
		if !ArrayCompare(actual, v, false) {
			t.Errorf("Split(%s) = %v; expected %v", k, actual, v)
		}
	}
}

func TestCleanTitle(t *testing.T) {
	cases := map[string][]string{
		"北区侦缉队 The Stronghold":   {"北区侦缉队", "The Stronghold"},
		"兴风作浪2 Trouble Makers":   {"兴风作浪2", "Trouble Makers"},
		"Tick Tick BOOM":         {"", "Tick Tick BOOM"},
		"比得兔2：逃跑计划":              {"比得兔2：逃跑计划", ""},
		"龙威山庄 99 Cycling Swords": {"龙威山庄", "99 Cycling Swords"},
		"我love你":                 {"我love你", ""},
		"我love 你":                {"我love 你", ""},
	}

	for title, want := range cases {
		chs, eng := SplitChsEngTitle(title)
		if chs != want[0] || eng != want[1] {
			t.Errorf("CleanTitle(%s) = %s-%s; want %s", title, chs, eng, want)
		}
	}
}

func TestCoverChsNumber(t *testing.T) {
	cases := map[string]int{
		"零":          0,
		"一":          1,
		"二":          2,
		"三":          3,
		"四":          4,
		"五":          5,
		"六":          6,
		"七":          7,
		"八":          8,
		"九":          9,
		"十":          10,
		"十一":         11,
		"十二":         12,
		"一十二":        12,
		"二十二":        22,
		"九十二":        92,
		"一百九十二":      192,
		"三千一百一十二":    3112,
		"三千一百九十二":    3192,
		"五万三千一百九十二":  53192,
		"五万零一百九十二":   50192,
		"五十三万零一百九十二": 530192,
		"五百万零一百九十二":  5000192,
		"四十二亿九千四百九十六万七千二百九十五": 4294967295,
	}
	for number, want := range cases {
		give := CoverChsNumber(number)
		if give != want {
			t.Errorf("CoverZhsNumber(%s) give %d, want %d", number, give, want)
		}
	}
}

func TestReplaceChsNumber(t *testing.T) {
	cases := map[string]string{
		"第一季":   "第1季",
		"第一集":   "第1集",
		"第十一集":  "第11集",
		"十一集":   "11集",
		"二":     "二",
		"一百九十二": "一百九十二",
	}
	for number, want := range cases {
		give := ReplaceChsNumber(number)
		if give != want {
			t.Errorf("ReplaceChsNumber(%s) give %s, want %s", number, give, want)
		}
	}
}

func TestSeasonCorrecting(t *testing.T) {
	cases := map[string]string{
		"邪恶力量.第01-14季.Supernatural.S01-S14.1080p.Blu-Ray.AC3.x265.10bit-Yumi": "邪恶力量.S01-S14.Supernatural.S01-S14.1080p.Blu-Ray.AC3.x265.10bit-Yumi",
		"堕落.第一季.2013.中英字幕￡CMCT无影":                                             "堕落.S01.2013.中英字幕￡CMCT无影",
		"一年一度喜剧大赛":                                                            "一年一度喜剧大赛",
		"亿万富犬.第一季":                                                            "亿万富犬.S01",
		"超级宝贝JOJO第二季":                                                         "超级宝贝JOJO.S02",
	}

	for title, want := range cases {
		give := SeasonCorrecting(title)
		if give != want {
			t.Errorf("SeasonCorrecting(%s) give: %s, want %s", title, give, want)
		}
	}
}

func TestEpisodeCorrecting(t *testing.T) {
	cases := map[string]string{
		"宝贝揪揪 第三季 第09集.mp4": "宝贝揪揪 第3季 E09.mp4",
		"宝贝揪揪 第三季 第01集.mp4": "宝贝揪揪 第3季 E01.mp4",
		"宝贝揪揪 第三季 第十集.mp4":  "宝贝揪揪 第3季 E10.mp4",
	}

	for title, want := range cases {
		give := EpisodeCorrecting(title)
		if give != want {
			t.Errorf("SeasonCorrecting(%s) give: %s, want %s", title, give, want)
		}
	}
}

func TestIsCollection(t *testing.T) {
	cases := map[string]bool{
		"邪恶力量.第01-14季.Supernatural.S01-S14.1080p.Blu-Ray.AC3.x265.10bit-Yumi": true,
		"外星也难民S01.Solar.Opposites.2020.1080p.WEB-DL.x265.AC3￡cXcY@FRDS":       false,
		"Heroes.S01-04.2006-2009.Complete.1080p.Amazon.Webdl.AVC.DDP5.1-DBTV": true,
	}

	for title, want := range cases {
		give := IsCollection(title)
		if give != want {
			t.Errorf("IsCollection(%s) give: %v, want %v", title, give, want)
		}
	}
}
