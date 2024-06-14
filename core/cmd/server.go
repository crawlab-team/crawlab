package cmd

import (
	"github.com/crawlab-team/crawlab/core/apps"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:     "server",
	Aliases: []string{"s"},
	Short:   "Start Crawlab server",
	Long:    `Start Crawlab node server that can serve as API, task scheduler, task runner, etc.`,
	Run: func(cmd *cobra.Command, args []string) {
		// app
		//svr := apps.GetServer(opts...)
		svr := apps.GetServerV2()

		// start
		apps.Start(svr)
	},
}
