package utils

import (
	"os"
	"pmimport/global"

	"go.uber.org/zap"
)

// @title    PathExists
// @description   文件目录是否存在
// @param     path            string
// @return    err             error

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// @title    createDir
// @description   批量创建文件夹
// @param     dirs            string
// @return    err             error

func CreateDir(dirs ...string) (err error) {
	for _, v := range dirs {
		exist, err := PathExists(v)
		if err != nil {
			return err
		}
		if !exist {
			err = os.MkdirAll(v, os.ModePerm)
			if err != nil {
				global.LOG.Debug("Create directory", zap.String("path", v), zap.Any(" error:", err))
			} else {
				global.LOG.Debug("Create directory", zap.String("path", v))
			}
		}
	}
	return err
}
