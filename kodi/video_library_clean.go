package kodi

import "fengqi/kodi-metadata-tmdb-cli/utils"

func (r *JsonRpc) AddCleanTask(directory string) {
	if !r.config.Enable {
		return
	}

	utils.Logger.DebugF("AddCleanTask %s", directory)
}
