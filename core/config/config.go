package config

import (
	"strings"

	"github.com/apex/log"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
	// Set default values
	c.setDefaults()

	// config
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // if config file is set, load it accordingly
	} else {
		viper.AddConfigPath("./conf") // if no config file is set, load by default
		viper.SetConfigName("config")
	}

	// config type as yaml
	viper.SetConfigType("yaml")

	// auto env
	viper.AutomaticEnv()

	// env prefix
	viper.SetEnvPrefix("CRAWLAB")

	// replacer
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// read in config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warn("No config file found. Using default values.")
		} else {
			log.Errorf("Error reading config file: %s", err)
			return err
		}
	}

	return nil
}

func (c *Config) setDefaults() {
	viper.SetDefault("edition", "global.edition.community")

	viper.SetDefault("mongo.host", "localhost")
	viper.SetDefault("mongo.port", 27017)
	viper.SetDefault("mongo.db", "crawlab_test")
	viper.SetDefault("mongo.username", "")
	viper.SetDefault("mongo.password", "")
	viper.SetDefault("mongo.authSource", "admin")

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8000)

	viper.SetDefault("grpc.address", "localhost:9666")
	viper.SetDefault("grpc.server.address", "0.0.0.0:9666")
	viper.SetDefault("grpc.authKey", "Crawlab2021!")

	viper.SetDefault("api.endpoint", "http://localhost:8000")

	viper.SetDefault("log.path", "/var/log/crawlab")
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
