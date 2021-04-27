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
	Short: "Start task handler",
	Long: `Start task handler service (worker) of Crawlab 
which runs tasks assigned by master node`,
	Run: func(cmd *cobra.Command, args []string) {
		handler := apps.NewHandler()
		handler.Init()
		handler.Run()
	},
}
