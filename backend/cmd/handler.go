package cmd

import (
	"crawlab/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(handlerCmd)
}

var handlerCmd = &cobra.Command{
	Use:   "handler",
	Short: "Start API server",
	Long:  `Start API server of Crawlab which serves data to frontend`,
	Run: func(cmd *cobra.Command, args []string) {
		handler := apps.NewHandler()
		handler.Init()
		handler.Run()
	},
}
