package cmd

import (
	"crawlab/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(schedulerCmd)
}

var schedulerCmd = &cobra.Command{
	Use:     "scheduler",
	Aliases: []string{"S"},
	Short:   "Start scheduler",
	Long: `Start a scheduler instance of Crawlab 
which assigns tasks to worker nodes to execute`,
	Run: func(cmd *cobra.Command, args []string) {
		scheduler := apps.NewScheduler()
		apps.Start(scheduler)
	},
}
