package cmd

import (
	"crawlab/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start API server",
	Long:  `Start API server of Crawlab which serves data to frontend`,
	Run: func(cmd *cobra.Command, args []string) {
		api := apps.NewApi()
		api.Init()
		api.Run()
	},
}
