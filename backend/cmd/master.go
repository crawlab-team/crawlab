package cmd

import (
	"crawlab/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(masterCmd)
}

var masterCmd = &cobra.Command{
	Use:     "master",
	Aliases: []string{"M"},
	Short:   "Start master",
	Long: `Start a master instance of Crawlab 
which runs api and assign tasks to worker nodes`,
	Run: func(cmd *cobra.Command, args []string) {
		// options
		var opts []apps.MasterOption

		// app
		master := apps.NewMaster(opts...)

		// start
		apps.Start(master)
	},
}
