package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"strconv"
	"strings"
)

func (d *Dir) ReadTvId() {
	idFile := d.Dir + "/" + d.OriginTitle + "/tmdb/id.txt"
	if _, err := os.Stat(idFile); err == nil {
		bytes, err := os.ReadFile(idFile)
		if err == nil {
			d.TvId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read tv id specially file: %s err: %v", idFile, err)
		}
	}
}

func (d *Dir) ReadSeason() {
	seasonFile := d.Dir + "/" + d.OriginTitle + "/tmdb/season.txt"
	if _, err := os.Stat(seasonFile); err == nil {
		bytes, err := os.ReadFile(seasonFile)
		if err == nil {
			d.Season, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
		} else {
			utils.Logger.WarningF("read season specially file: %s err: %v", seasonFile, err)
		}
	}

	if d.Season == 0 && len(d.YearRange) == 0 {
		d.Season = 1
	}
}

func (d *Dir) ReadGroupId() {
	groupFile := d.Dir + "/" + d.OriginTitle + "/tmdb/group.txt"
	if _, err := os.Stat(groupFile); err == nil {
		bytes, err := os.ReadFile(groupFile)
		if err == nil {
			d.GroupId = strings.Trim(string(bytes), "\r\n ")
		} else {
			utils.Logger.WarningF("read group id specially file: %s err: %v", groupFile, err)
		}
	}
}
