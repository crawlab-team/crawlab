package config

import (
	"github.com/crawlab-team/crawlab/core/constants"
	"github.com/crawlab-team/crawlab/core/entity"
	"github.com/crawlab-team/crawlab/core/utils"
	"github.com/spf13/viper"
)

type Config entity.NodeInfo

type Options struct {
	Key        string
	Name       string
	IsMaster   bool
	AuthKey    string
	MaxRunners int
}

var DefaultMaxRunner = 8

var DefaultConfigOptions = &Options{
	Key:        utils.NewUUIDString(),
	IsMaster:   utils.IsMaster(),
	AuthKey:    constants.DefaultGrpcAuthKey,
	MaxRunners: 0,
}

func NewConfig(opts *Options) (cfg *Config) {
	if opts == nil {
		opts = DefaultConfigOptions
	}
	if opts.Key == "" {
		if viper.GetString("node.key") != "" {
			opts.Key = viper.GetString("node.key")
		} else {
			opts.Key = utils.NewUUIDString()
		}
	}
	if opts.Name == "" {
		if viper.GetString("node.name") != "" {
			opts.Name = viper.GetString("node.name")
		} else {
			opts.Name = opts.Key
		}
	}
	if opts.AuthKey == "" {
		if viper.GetString("grpc.authKey") != "" {
			opts.AuthKey = viper.GetString("grpc.authKey")
		} else {
			opts.AuthKey = constants.DefaultGrpcAuthKey
		}
	}
	if opts.MaxRunners == 0 {
		if viper.GetInt("task.handler.maxRunners") != 0 {
			opts.MaxRunners = viper.GetInt("task.handler.maxRunners")
		} else {
			opts.MaxRunners = DefaultMaxRunner
		}
	}
	return &Config{
		Key:        opts.Key,
		Name:       opts.Name,
		IsMaster:   opts.IsMaster,
		AuthKey:    opts.AuthKey,
		MaxRunners: opts.MaxRunners,
	}
}
