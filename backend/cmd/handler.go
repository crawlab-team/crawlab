package cmd

import (
	"crawlab/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(handlerCmd)
}

var handlerCmd = &cobra.Command{
	Use:     "handler",
	Aliases: []string{"H"},
	Short:   "Start handler",
	Long: `Start a handler instance of Crawlab 
which runs tasks with given parameters`,
	Run: func(cmd *cobra.Command, args []string) {
		handler := apps.NewHandler()
		apps.Start(handler)
	},
}
