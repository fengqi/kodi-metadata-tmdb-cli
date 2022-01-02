package utils

import (
	"testing"
)

func TestCleanTitle(t *testing.T) {
	unit := map[string]string{
		"中国通史":                           "中国通史",
		"The Sex Lives of College":       "The Sex Lives of College",
		"BBC 行星 The Planets":             "The Planets",
		"[Arcane]":                       "Arcane",
		"[机智的医生生活]":                      "机智的医生生活",
		"[机智的医生生活] A Wise Doctor's Life": "A Wise Doctor's Life",
	}

	for k, v := range unit {
		actual := CleanTitle(k)
		if actual != v {
			t.Errorf("cleanTitle(%s) = %s; expected %s", k, actual, v)
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
