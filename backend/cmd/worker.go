package cmd

import (
	"crawlab/apps"
	"fmt"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	workerConfigPath  string
	workerGrpcAddress string
	workerGrpcAuthKey string
)

func init() {
	rootCmd.AddCommand(workerCmd)

	workerCmd.PersistentFlags().StringVarP(&workerConfigPath, "config-path", "c", "", "Config path of worker node")
	_ = viper.BindPFlag("configPath", workerCmd.PersistentFlags().Lookup("configPath"))

	workerCmd.PersistentFlags().StringVarP(&workerGrpcAddress, "grpc-address", "g", "", "gRPC address of worker node")
	_ = viper.BindPFlag("grpcAddress", workerCmd.PersistentFlags().Lookup("grpcAddress"))

	workerCmd.PersistentFlags().StringVarP(&workerGrpcAuthKey, "grpc-auth-key", "a", "", "gRPC auth key of worker node")
	_ = viper.BindPFlag("grpcAuthKey", workerCmd.PersistentFlags().Lookup("grpcAuthKey"))
}

var workerCmd = &cobra.Command{
	Use:     "worker",
	Aliases: []string{"W"},
	Short:   "Start worker",
	Long: `Start a worker instance of Crawlab 
serving in the worker node and executes tasks
assigned by the master node`,
	Run: func(cmd *cobra.Command, args []string) {
		// options
		var opts []apps.WorkerOption
		if workerConfigPath != "" {
			opts = append(opts, apps.WithWorkerConfigPath(workerConfigPath))
			viper.Set("config.path", workerConfigPath)
		}
		if workerGrpcAddress != "" {
			address, err := entity.NewAddressFromString(workerGrpcAddress)
			if err != nil {
				fmt.Println(fmt.Sprintf("invalid grpc-address: %s", workerGrpcAddress))
				return
			}
			opts = append(opts, apps.WithWorkerGrpcAddress(address))
			viper.Set("grpc.address", workerGrpcAddress)
		}
		if workerGrpcAuthKey != "" {
			viper.Set("grpc.authKey", workerGrpcAuthKey)
		}

		// app
		master := apps.NewWorker(opts...)

		// start
		apps.Start(master)
	},
}
