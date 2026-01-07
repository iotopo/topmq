package cmd

import (
	"fmt"
	"github.com/iotopo/topmq/config"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print the version and release date",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s build at %s\n", config.Version, config.ReleaseTime)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and release date",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("v%s build at %s\n", config.Version, config.ReleaseTime)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(versionCmd)
}
