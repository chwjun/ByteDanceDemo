// Package cmd @Author: youngalone [2023/8/7]
package cmd

import (
	"bytedancedemo/cmd/api"
	"bytedancedemo/cmd/migrate"
	"errors"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var RootCmd = &cobra.Command{
	Use:               "ByteDance",
	SilenceUsage:      true,
	DisableAutoGenTag: false,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("欢迎使用抖音极简版 -h查看命令")
	},
}

func init() {
	RootCmd.AddCommand(migrate.StartCmd)
	RootCmd.AddCommand(api.StartCmd)
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
