package cmd

import (
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile     string
	runOnMaster bool

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

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVar(&runOnMaster, "runOnMaster", false, "Whether to run tasks on master node (default: false)")
	_ = viper.BindPFlag("runOnMaster", rootCmd.PersistentFlags().Lookup("runOnMaster"))
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
