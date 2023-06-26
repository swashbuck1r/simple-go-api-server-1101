package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/swashbuck1r/simple-go-api-server/server"
)

func init() {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start server",
		Long:  `Start the api-server`,
		Run: func(cmd *cobra.Command, args []string) {
			var config server.Configuration
			viper.Unmarshal(&config)
			server.Start(&config)
		},
	}

	// port flag
	cmd.Flags().IntP("port", "p", viper.GetInt("port"), "Port to run the server on")
	viper.BindPFlag("port", cmd.Flags().Lookup("port"))

	rootCmd.AddCommand(cmd)
}
