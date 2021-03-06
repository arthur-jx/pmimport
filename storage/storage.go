package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"pmimport/global"
	"pmimport/utils"
	"strings"
	"time"
)

//文件存储API接口， 方便以后扩展为别的存储系统

//返回指定的文件是否存在，或出错
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func PathIsDir(path string) bool {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir()
	}
	return false
}

//给定路径，创建子目录
// func CreateMediaDir(basedir string, t time.Time) bool {
// 	if os.Chdir(basedir) == nil {
// 		dir := fmt.Sprintf("%d/%s", t.Year(), t.Local().Format("2006-01-02"))
// 		if os.MkdirAll(dir, os.ModeDir|os.ModeSetuid|os.ModeSetgid) == nil {
// 			return true
// 		}
// 	}
// 	return false
// }

func CopyFile(src, dest, fileHash string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destDir := filepath.Dir(dest)
	if ok, _ := utils.PathExists(destDir); !ok {
		utils.CreateDir(destDir)
	}

	destination, err := os.Create(dest)
	if err != nil {
		return err
	}

	size, err := io.Copy(destination, source)
	destination.Close()

	if size != sourceFileStat.Size() {
		err = fmt.Errorf("copy file fail.")
	} else {
		//verify target file Hash
		newHash, e := utils.GetFileSHA1(dest)
		if e == nil {
			if strings.Compare(newHash, fileHash) != 0 {
				err = fmt.Errorf("copy file verify fail.")
			}
		} else {
			err = fmt.Errorf("copy file verify fail.")
		}
	}

	return err
}

//返回媒体文件的保存目录
func GetImportStoragePath(createTime time.Time) string {
	userPath := GetUserStoragePath()

	if len(userPath) == 0 {
		return ""
	}

	//path: storage/user/media/year/date
	path := fmt.Sprintf("%s/media/%d/%s", userPath, createTime.Year(), createTime.Local().Format("2006-01-02"))
	return path
}

func GetUserStoragePath() string {
	if len(global.CONFIG.Storage.UserId) == 0 {
		return ""
	}

	return fmt.Sprintf("%s/%s", global.CONFIG.Storage.Path, global.CONFIG.Storage.UserId)
}

func GetUserMediaFilePath() string {
	path := fmt.Sprintf("%s/media", GetUserStoragePath())

	if ok, _ := utils.PathExists(path); !ok {
		dirs := []string{path}
		err := utils.CreateDir(dirs...)
		if err != nil {
			fmt.Printf("Open log file error:%s", err.Error())
			os.Exit(1)
		}
	}

	return path
}

func GetUserLogsFilePath() string {
	path := fmt.Sprintf("%s/logs", GetUserStoragePath())

	if ok, _ := utils.PathExists(path); !ok {
		return ""
	}

	path += fmt.Sprintf("/import_%s.log", time.Now().Local().Format("20060102_150405"))
	return path
}

func RenameFile(src, dest string) error {
	return os.Rename(src, dest)
}
