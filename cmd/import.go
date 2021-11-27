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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"pmimport/global"
	"pmimport/media"
	"pmimport/utils"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type ImportArgs struct {
	importFrom    string //要导入的源，可以是一个目录，以可以是一个媒体文件
	overwrite     bool   //如果文件相同是否覆盖, 优先于rename
	backup        bool   //如果目标文件存在，是否改名导入
	interactive   bool   //覆盖前是否提示
	recursive     bool   //是否导出子目录文件
	destroy       bool   //是否删除已导入的源文件, 优先于change参数
	rename        bool   //是否对已导入的文件重命令，使用固定后缀
	useCreateDate bool   //允许导入在Exif信息中没有拍摄日期的照片，并使用创建日期作为归档日期,
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// importCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	importCmd.Flags().StringVar(&importArgs.importFrom, "from", "./", "import media source path")
	rootCmd.MarkFlagRequired("src")

	importCmd.Flags().BoolVar(&importArgs.overwrite, "overwrite", false, "remove of each existing destination file (default:true)")
	importCmd.Flags().BoolVar(&importArgs.backup, "backup", false, "make a backup of each existing destination file")
	importCmd.Flags().BoolVar(&importArgs.interactive, "interactive", false, "prompt before overwrite")
	importCmd.Flags().BoolVar(&importArgs.recursive, "recursive", false, "import directories recursively")
	importCmd.Flags().BoolVar(&importArgs.destroy, "destroy", false, "destroy source file of import success")
	importCmd.Flags().BoolVar(&importArgs.rename, "rename", false, "rename source file of import success")

	importCmd.Flags().BoolVar(&importArgs.useCreateDate, "use-create-date", false, "use create date import Not Exif Info photo")
	importCmd.Flags().StringVar(&importArgs.useModel, "model", "", "setting import default model")
	importCmd.Flags().BoolVar(&importArgs.timeIsLocal, "islocaltime", false, "media time is local")

	importCmd.Flags().StringVar(&importArgs.tags, "tags", "", "tags for media files")
	importCmd.MarkFlagRequired("tags")

	importCmd.Flags().BoolVar(&importArgs.verbose, "verbose", false, "explain what is being done")

}

func importCommand(cmd *cobra.Command, args []string) {
	fmt.Println("import called")
	if len(importArgs.importFrom) == 0 {
		return
	}
	srcPath, err := filepath.Abs(importArgs.importFrom)
	if err == nil {
		fmt.Println("FORM:", srcPath)
		importPath(srcPath, &importArgs)
	}
}

func importPath(srcPath string, args *ImportArgs) (err error) {
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
						importPath(nextName, args)
					}
				} else {
					ImportFile(nextName, args)
				}
			}
		}
	} else {
		ImportFile(srcPath, args)
	}

	return nil
}

func ImportFile(filePath string, args *ImportArgs) (err error) {
	fmt.Println("Import:", filePath)
	createTime, tags, err := getFileTimeAndTags(filePath)
	if err == nil {
		fmt.Printf("   -> file info:   %v \t [%v]\n", createTime.Local().Format("2006-01-02 15:04:05"), tags)

		//TODO:: move or copy file to storage
	} else {
		fmt.Printf("   -> ERR:%v\n", err.Error())
	}

	fmt.Println("    FAIL")
	return nil
}

func getFileTimeAndTags(filePath string) (createTime time.Time, tags string, err error) {
	exif, e := media.GetExif(filePath)
	if e != nil {
		err = e
		return
	}

	//Get field: Model, CreateDate, OffsetTimeOriginal,LensModel
	Model, _ := media.GetExifInfoString(exif, "Model")

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

	LensModel, _ := media.GetExifInfoString(exif, "LensModel")

	if len(Model) == 0 {
		if len(importArgs.useModel) > 0 {
			Model = importArgs.useModel
		}
	}

	fmt.Printf("   -> Exif Info:[%s]\t[%s]\t[%s]\t[%s]\n", Model, CreateDate, OffsetTimeOriginal, LensModel)

	if len(Model) > 0 {
		tags += Model
		if len(LensModel) > 0 {
			tags += "," + LensModel
		}
		tags = strings.Replace(tags, " ", "_", -1)
	} else {
		err = fmt.Errorf("invalid photo model")
		return
	}

	if len(CreateDate) > 0 {
		layout := "2006:01:02 15:04:05"
		if len(OffsetTimeOriginal) > 0 {
			layout += "-07:00"
			CreateDate += OffsetTimeOriginal
		} else {
			if importArgs.timeIsLocal {
				layout += "-07:00"
				CreateDate += "+08:00"
			}
		}

		ctime, e := time.Parse(layout, CreateDate)
		if e == nil {
			createTime = ctime
			if len(tags) > 0 {
				err = nil
				return
			}
		}
	}
	err = fmt.Errorf("invalid photo time")

	return
}
