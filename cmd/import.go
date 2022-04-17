/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"pmimport/global"
	"pmimport/media"
	"pmimport/utils"
	"strings"
	"time"

	"pmimport/storage"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ImportArgs struct {
	importFrom    string //要导入的源，可以是一个目录，以可以是一个媒体文件
	overwrite     bool   //如果文件相同是否覆盖, 优先于rename
	nobackup      bool   //如果目标文件存在，是否改名
	interactive   bool   //覆盖前是否提示
	recursive     bool   //递归导出子目录文件
	destroy       bool   //是否删除已导入的源文件, 优先于change参数
	rename        bool   //是否对已导入的文件重命令，使用固定后缀
	moveto        string //将成功导入的源文件移动到指定目录
	excludeFile   string //不导入指定文件名前缀的文件
	setCreateTime string //导入在Exif信息中没有拍摄日期的照片时，使用指定日期作为归档日期,
	useModel      string //对于没有相机Model信息的，使用指定的model
	timeIsLocal   bool   //文件中的时间为本地地间
	tags          string //文件的标签
	verbose       bool   //是否显示导入过程祥情
	test          bool   //只执行导入过程，并显示文件处理，但是不实际复制文件
}

var importArgs ImportArgs

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: importCommand,
}

func init() {
	rootCmd.AddCommand(importCmd)

	importCmd.Flags().StringVar(&importArgs.importFrom, "from", "./", "import media source path")
	rootCmd.MarkFlagRequired("src")

	importCmd.Flags().BoolVar(&importArgs.overwrite, "overwrite", false, "remove of each existing destination file (default:true)")
	importCmd.Flags().BoolVar(&importArgs.nobackup, "nobackup", false, "don't make a backup of each overwrite existing destination file")
	importCmd.Flags().BoolVar(&importArgs.interactive, "interactive", false, "prompt before overwrite")
	importCmd.Flags().BoolVar(&importArgs.recursive, "recursive", false, "recursively import subdirectories")
	importCmd.Flags().BoolVar(&importArgs.destroy, "destroy", false, "destroy source file of import success")
	importCmd.Flags().BoolVar(&importArgs.rename, "rename", false, "rename source file of import success")
	importCmd.Flags().StringVar(&importArgs.moveto, "moveto", "", "move to the path of source file import success")

	importCmd.Flags().StringVar(&importArgs.excludeFile, "exclude-file", "", "exclude file with spacified filename prefix")

	importCmd.Flags().StringVar(&importArgs.setCreateTime, "set-create-time", "", "set time with import file, format: YYYY-MM-DDTmm:hh:ss")
	importCmd.Flags().StringVar(&importArgs.useModel, "model", "", "setting import default model")
	importCmd.Flags().BoolVar(&importArgs.timeIsLocal, "islocaltime", true, "media time is local")

	importCmd.Flags().StringVar(&importArgs.tags, "tags", "", "tags for media files")
	importCmd.MarkFlagRequired("tags")

	importCmd.Flags().BoolVar(&importArgs.verbose, "verbose", false, "explain what is being done")

}

func importCommand(cmd *cobra.Command, args []string) {
	if len(importArgs.importFrom) == 0 {
		return
	}

	//TODO:: 检查用户media 目录必须存在
	if !storage.FileExist(storage.GetUserMediaFilePath()) {
		global.LOG.Error("User media dir don't exist.")
		os.Exit(-1)
	}

	srcPath, err := filepath.Abs(importArgs.importFrom)
	if err == nil {
		global.LOG.Info("IMPORT TIME:", zap.Time("time", time.Now()))
		global.LOG.Info("FORM:", zap.String("path:", srcPath))
		global.LOG.Info("User:", zap.String("id", global.CONFIG.Storage.UserId))
		importPath(srcPath, srcPath, &importArgs)
	}
}

func importPath(srcPath string, fromPath string, args *ImportArgs) (err error) {
	if ok, e := utils.PathExists(srcPath); !ok {
		err = e
		global.LOG.Error("Source media path is don't exists", zap.String("path", srcPath))
		return
	}

	fi, err := os.Lstat(srcPath)
	if fi.IsDir() {
		//遍历目录
		entrys, err := ioutil.ReadDir(srcPath)
		if err == nil {
			for _, info := range entrys {
				nextName := srcPath + string(os.PathSeparator) + info.Name()
				if info.IsDir() {
					if args.recursive {
						importPath(nextName, fromPath, args)
					}
				} else {
					srcFileFull := strings.Replace(nextName, fromPath, "", 1)
					ImportFile(nextName, srcFileFull, args)
				}
			}
		}
	} else {
		srcFileFull := strings.Replace(srcPath, fromPath, "", 1)
		if len(srcFileFull) == 0 {
			srcFileFull = path.Base(srcPath)
		}
		ImportFile(srcPath, srcFileFull, args)
	}

	return nil
}

//filePath: 导入文件的全路径
//srcFileFull: 导入文件不含ＦＯＲＭ根目录的路径
func ImportFile(srcFile string, srcFileFull string, args *ImportArgs) (err error) {
	global.LOG.Debug("Import", zap.String("path", srcFile))

	if len(args.excludeFile) > 0 {
		fileName := path.Base(srcFile)
		if strings.HasPrefix(fileName, args.excludeFile) {
			global.LOG.Info("[SKIP]", zap.String("src", srcFileFull))
			return nil
		}
	}

	mediaInfo, err := getFileInfos(srcFile, args)
	if err == nil {
		global.LOG.Debug("file info", zap.Any("time", mediaInfo.CreateTime.Local().Format("2006-01-02 15:04:05")), zap.Any("model", mediaInfo.Model))

		mediaInfo.Tags = strings.Join([]string{mediaInfo.Tags, args.tags}, global.TagsSplit)
		mod := "[NO]"
		srcExtMode := ""
		destExtMode := ""

		importPath := storage.GetImportStoragePath(mediaInfo.CreateTime)
		targetPath := ""
		targetName := "" //不含库路径的导入路径文件名
		//get file sha256
		filesha, e := utils.GetFileSHA1(srcFile)
		if e == nil {
			mediaInfo.FileHash = filesha

			targetPath = path.Join(importPath, filesha+path.Ext(srcFile))
			targetName = strings.Replace(targetPath, storage.GetUserStoragePath(), "", 1)

			hasCopy := true //是否需要复制源文件
			if storage.FileExist(targetPath) {
				if importArgs.overwrite {
					//文件已存在
					if !importArgs.nobackup {
						//备份目标文件
						newName := strings.Replace(path.Base(targetPath), path.Ext(targetPath), "", 1) +
							time.Now().Format("_20060102150405999") + path.Ext(targetPath)
						err := storage.RenameFile(targetPath, path.Join(path.Dir(targetPath), newName))
						global.LOG.Debug("[RENAME]", zap.String("old name", targetName), zap.String("new name", newName), zap.Error(err))
						destExtMode = "BAK"
					} else {
						//覆盖目标文件
						destExtMode = "OVER"
					}
				} else {
					destExtMode = "NO"
					// global.LOG.Info("[NO]", zap.String("src", srcFileFull), zap.String("dest", targetName), zap.Any("error", "dest file is exist."))
					// return fmt.Errorf("media file exist in storage")
				}
			}

			isFinish := false //文件是否已复制

			if hasCopy {
				//copy file
				err := storage.CopyFile(srcFile, targetPath, filesha)
				if err != nil {
					global.LOG.Error("[COPY]", zap.String("src", srcFileFull), zap.String("dest", targetName), zap.String("Error", err.Error()))
					return err
				} else {
					mod = "[FIN]"
					isFinish = true
				}
			} else {
				isFinish = true
			}

			if isFinish {
				//save or update fileinfo
				updateMediaInfoFiles(targetPath, mediaInfo)

				//remove OR rename source file
				if importArgs.destroy {
					//remove source modia file
					err := os.Remove(srcFile)
					srcExtMode = "RM"
					if err != nil {
						srcExtMode = "RM_ERR"
					}
				} else {
					if importArgs.rename {
						//rename source file
						newName := "import-" + strings.Replace(path.Base(srcFile), path.Ext(srcFile), "", 1) + path.Ext(srcFile)
						err := storage.RenameFile(srcFile, path.Join(path.Dir(srcFile), newName))
						if err != nil {
							srcExtMode = "RENAME_ERR"
						} else {
							srcExtMode = "RENAME"
						}
					}
				}
			}

			global.LOG.Info(mod, zap.String("src", srcExtMode+">"+srcFileFull), zap.String("dest", destExtMode+">"+targetName))
		}
	} else {
		global.LOG.Info("[NO]", zap.String("src", srcFileFull), zap.String("error", err.Error()))
	}

	return nil
}

//保存媒体信息文件，目录已经有文件了，就合并文件中对应的字段
func updateMediaInfoFiles(mediaPath string, info *media.MediaFileInfo) {
	fileName := path.Base(mediaPath)
	fileName = strings.Replace(fileName, path.Ext(mediaPath), "", 1)

	filePath := path.Join(path.Dir(mediaPath), fileName+"_info.json")

	if storage.FileExist(filePath) {
		//load old info
		f, err := ioutil.ReadFile(filePath)
		if err == nil {
			var oldInfo media.MediaFileInfo
			err = json.Unmarshal(f, &oldInfo)
			if err == nil {
				if len(oldInfo.LensModel) > 0 {
					if len(info.LensModel) == 0 {
						info.LensModel = oldInfo.LensModel
					}
				}

				if len(oldInfo.LatLong) > 0 {
					if len(info.LatLong) == 0 {
						info.LatLong = oldInfo.LatLong
					}
				}

				strings.Join([]string{oldInfo.Tags, info.Tags}, global.TagsSplit)

				info.AlbumText = oldInfo.AlbumText + info.AlbumText
				info.Remark = oldInfo.Remark + info.Remark
			}
		}
	}

	infoBuff, err := json.Marshal(info)
	if err != nil {
		global.LOG.Error("Save media info file error", zap.String("file", filePath), zap.Error(err))
	} else {
		err = ioutil.WriteFile(filePath, infoBuff, 0666)
		if err != nil {
			global.LOG.Error("save media info file error", zap.String("file", filePath), zap.Error(err))
		}
	}
}

func getFileInfos(filePath string, args *ImportArgs) (info *media.MediaFileInfo, err error) {
	exif, e := media.GetExif(filePath)
	if e != nil {
		err = e
		return
	}

	var fileInfo media.MediaFileInfo

	Model, _ := media.GetExifInfoString(exif, "Model")
	if len(Model) == 0 {
		if len(importArgs.useModel) > 0 {
			Model = importArgs.useModel
		}
	}

	LensModel, _ := media.GetExifInfoString(exif, "LensModel")

	if len(Model) == 0 {
		err = fmt.Errorf("invalid photo model")
		return
	}

	//Get field: Model, CreateDate, OffsetTimeOriginal,LensModel
	var fileTime time.Time
	if len(args.setCreateTime) > 0 {
		if "now" == args.setCreateTime {
			fileTime = time.Now()
		} else {
			var err error
			loc, _ := time.LoadLocation("Local")
			fileTime, err = time.ParseInLocation("2006-1-2T15:4:5", args.setCreateTime, loc)
			if err != nil {
				return nil, fmt.Errorf("set create date invalid:%s", args.setCreateTime)
			}
		}
	} else {
		OffsetTimeOriginal, _ := media.GetExifInfoString(exif, "OffsetTimeOriginal")
		//视频文件读取CreationDate
		CreateDate, ok := media.GetExifInfoString(exif, "CreationDate")
		if !ok {
			CreateDate, _ = media.GetExifInfoString(exif, "CreateDate")
		}
		//提取时区标记
		pos := strings.Index(CreateDate, "+")
		if pos > 0 {
			if len(OffsetTimeOriginal) == 0 {
				OffsetTimeOriginal = CreateDate[pos:]
			}
			CreateDate = CreateDate[0:pos]
		}

		if len(CreateDate) > 0 {
			useLocTime := false
			layout := "2006:01:02 15:04:05"
			if len(OffsetTimeOriginal) > 0 {
				layout += "-07:00"
				CreateDate += OffsetTimeOriginal
			} else {
				if importArgs.timeIsLocal {
					useLocTime = true
					// layout += "-07:00"
					// CreateDate += "+08:00"
				}
			}

			global.LOG.Debug("Exif Info", zap.String("Model", Model), zap.String("CreateDate", CreateDate),
				zap.String("Timezone", OffsetTimeOriginal), zap.String("LensModel", LensModel))

			var ctime time.Time
			var e error
			if useLocTime {
				loc := time.Now().Location()
				ctime, e = time.ParseInLocation(layout, CreateDate, loc)
			} else {
				ctime, e = time.Parse(layout, CreateDate)
			}
			if e == nil {
				fileTime = ctime
			} else {
				return nil, fmt.Errorf("invalid photo time")
			}
		}
	}

	if !fileTime.IsZero() {
		//信息至少应该有日期和相机型号
		if len(Model) > 0 {
			fileInfo.CreateTime = fileTime
			fileInfo.Model = Model
			fileInfo.LensModel = LensModel

			latLong, _ := media.GetExifLanLong(exif)

			fileInfo.LatLong = latLong

			info = &fileInfo
			err = nil
			return
		} else {
			return nil, fmt.Errorf("invalid photo model")
		}
	} else {
		return nil, fmt.Errorf("invalid photo time")
	}
}
