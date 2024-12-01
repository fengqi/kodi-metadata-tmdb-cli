package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/fs"
	"os"
)

func listFile(dir string) []fs.DirEntry {
	entry, err := os.ReadDir(dir)
	if err != nil {
		utils.Logger.ErrorF("list dir: %s err: %v", dir, err)
		return nil
	}

	items := make([]fs.DirEntry, 0)
	for _, item := range entry {
		// 过滤无用文件
		if item.Name()[0:1] == "." || utils.InArray(config.Collector.SkipFolders, item.Name()) {
			utils.Logger.DebugF("pass file: %s", item.Name())
			return nil
		}

		items = append(items, item)
	}

	return items
}
