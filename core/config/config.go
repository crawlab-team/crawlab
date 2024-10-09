package config

import (
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab/trace"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"strings"
)

func init() {
	// config instance
	c := Config{Name: ""}

	// init config file
	if err := c.Init(); err != nil {
		log.Warn("unable to init config")
		return
	}

	// watch config change and load responsively
	c.WatchConfig()

	// init log level
	c.initLogLevel()
}

type Config struct {
	Name string
}

type InitConfigOptions struct {
	Name string
}

func (c *Config) WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

func (c *Config) Init() (err error) {
	// config
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // if config file is set, load it accordingly
	} else {
		viper.AddConfigPath("./conf") // if no config file is set, load by default
		viper.SetConfigName("config")
	}

	// config type as yaml
	viper.SetConfigType("yaml") // default yaml

	// auto env
	viper.AutomaticEnv() // load matched environment variables

	// env prefix
	viper.SetEnvPrefix("CRAWLAB") // environment variable prefix as CRAWLAB

	// replacer
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// read in config
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("Error reading config file, %s", err)
		trace.PrintError(err)
		return err
	}

	return nil
}

func (c *Config) initLogLevel() {
	// set log level
	logLevel := viper.GetString("log.level")
	l, err := log.ParseLevel(logLevel)
	if err != nil {
		l = log.InfoLevel
	}
	log.SetLevel(l)
}
