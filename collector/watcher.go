package collector

import (
	"os"
)

// watcherCallback 监听文件变化的回调函数
func (c *collector) watcherCallback(filename string, fileInfo os.FileInfo) {
	//	log.Println(filename, fileInfo)
}
