package cmd

import (
	"fmt"
	"github.com/iotopo/topmq/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"os"
)

var rootCmd = &cobra.Command{
	Use:          config.AppName,
	Version:      config.Version,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func AddCommand(cmds ...*cobra.Command) {
	rootCmd.AddCommand(cmds...)
}

func Execute() {
	logrus.SetLevel(logrus.ErrorLevel)
	logrus.SetOutput(io.Discard)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func init() {
	// 隐藏根命令的用法信息
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
}
