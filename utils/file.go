package utils

import (
	"fmt"
	"io"
	"os"
)

func IsDir(dir string) bool {
	return false
}

func DirExist(dir string) bool {
	return false
}

func IsFile(file string) bool {
	return false
}

func FileExist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}

	return false
}

func CopyFile(dstName, srcName string) (writeen int64, err error) {
	src, err := os.Open(dstName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()

	dst, err := os.OpenFile(srcName, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}
