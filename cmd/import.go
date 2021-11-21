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
	tags          string //文件的标签
	verbose       bool   //是否显示导入过程祥情
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
	}

	return nil
}

func ImportFile(filePath string, args *ImportArgs) (err error) {
	fmt.Println("Import:", filePath)
	//TODO:: read Exif data
	exif := media.GetExif(filePath)
	media.ShowExit(exif)

	fmt.Println("    FAIL")
	return nil
}
