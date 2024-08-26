package kodi

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
)

func (r *JsonRpc) AddCleanTask(directory string) {
	if !config.Kodi.Enable {
		return
	}

	utils.Logger.DebugF("AddCleanTask %s", directory)
}
