package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "crawlab",
		Short: "CLI tool for Crawlab",
		Long: `The CLI tool is for controlling against Crawlab.
Crawlab is a distributed web crawler and task admin platform
aimed at making web crawling and task management easier.
`,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

// GetRootCmd get rootCmd instance
func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "c", "", "Use Custom Config File")
}
