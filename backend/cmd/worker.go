package cmd

import (
	"crawlab/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(workerCmd)
}

var workerCmd = &cobra.Command{
	Use:     "worker",
	Aliases: []string{"W"},
	Short:   "Start worker",
	Long: `Start a worker instance of Crawlab 
serving in the worker node and executes tasks
assigned by the master node`,
	Run: func(cmd *cobra.Command, args []string) {
		worker := apps.NewWorker()
		apps.Start(worker)
	},
}
