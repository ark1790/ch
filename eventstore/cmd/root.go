package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// RootCmd is the root command of the server
	RootCmd = &cobra.Command{
		Use:   "eventstore",
		Short: "eventstore is a gRPC API server",
		Long:  "eventstore is a gRPC API server",
	}
)

func init() {
	viper.AutomaticEnv()
}
