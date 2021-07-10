package cmd

import (
	"crawlab/apps"
	"fmt"
	"github.com/crawlab-team/crawlab-core/entity"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	runOnMaster       bool
	masterConfigPath  string
	masterGrpcAddress string
)

func init() {
	rootCmd.AddCommand(masterCmd)

	masterCmd.PersistentFlags().StringVarP(&masterConfigPath, "config-path", "c", "", "Config path of master node")
	_ = viper.BindPFlag("configPath", masterCmd.PersistentFlags().Lookup("configPath"))

	masterCmd.PersistentFlags().StringVarP(&masterGrpcAddress, "grpc-address", "g", "", "gRPC address of master node")
	_ = viper.BindPFlag("grpcAddress", masterCmd.PersistentFlags().Lookup("grpcAddress"))
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
		if masterConfigPath != "" {
			opts = append(opts, apps.WithMasterConfigPath(masterConfigPath))
			viper.Set("config.path", masterConfigPath)
		}
		opts = append(opts, apps.WithRunOnMaster(runOnMaster))
		if masterGrpcAddress != "" {
			address, err := entity.NewAddressFromString(masterGrpcAddress)
			if err != nil {
				fmt.Println(fmt.Sprintf("invalid grpc-address: %s", masterGrpcAddress))
			}
			opts = append(opts, apps.WithMasterGrpcAddress(address))
			viper.Set("grpc.address", masterGrpcAddress)
			viper.Set("grpc.server.address", masterGrpcAddress)
		}

		// app
		master := apps.NewMaster(opts...)

		// start
		apps.Start(master)
	},
}
