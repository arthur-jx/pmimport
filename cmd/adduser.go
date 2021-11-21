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
	"os"
	"pmimport/global"
	"pmimport/utils"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// adduserCmd represents the adduser command
var adduserCmd = &cobra.Command{
	Use:   "adduser",
	Short: "add a user dir in to storage",
	Long:  ``,
	Run:   addCommand,
}

func init() {
	rootCmd.AddCommand(adduserCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adduserCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adduserCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func addCommand(cmd *cobra.Command, args []string) {
	storage := global.CONFIG.Storage.Path
	userid := global.CONFIG.Storage.UserId

	userPath := storage + string(os.PathSeparator) + userid
	if ok, _ := utils.PathExists(userPath); ok {
		global.LOG.Error("Storage user Id is exists, don't create", zap.String("userid", userid))
		return
	}

	err := utils.CreateDir(userPath)
	if err != nil {
		global.LOG.Error("Create user dir error", zap.Any("error", err))
		return
	}

	dirs := []string{userPath + string(os.PathSeparator) + "media",
		userPath + string(os.PathSeparator) + "photo_album"}
	utils.CreateDir(dirs...)

	fmt.Println("User create SUCCESS")
}
