package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("./conf")
		viper.SetConfigName("config")
	}

	// file format as yaml
	viper.SetConfigType("yaml")

	// auto load env
	viper.AutomaticEnv()

	// env prefix as CRAWLAB
	viper.SetEnvPrefix("CRAWLAB")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
