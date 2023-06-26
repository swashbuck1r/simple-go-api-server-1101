package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/swashbuck1r/simple-go-api-server/config"
)

var (
	// Used for flags.
	cfgFile string
	envMode string

	rootCmd = &cobra.Command{
		Use:   "api-server",
		Short: "api-server commands",
		Long:  `Commands to operate the api-server.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()

			var config config.Configuration
			viper.Unmarshal(&config)
			fmt.Println("\n[server configuration]\n", PrettyPrint(config))
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// initialize default configuration
	viper.SetConfigType("yaml")
	if err := viper.ReadConfig(bytes.NewBuffer(config.DefaultConfiguration)); err != nil {
		log.Fatalln("Failed to read default config file", err)
	}

	// config file flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// env mode flag
	rootCmd.PersistentFlags().StringVar(&envMode, "env", "dev", "environment execution mode (dev, production, etc)")
	viper.BindPFlag("env", rootCmd.Flags().Lookup("env"))

	// initialize the viper configuration based on the config file and env mode
	cobra.OnInitialize(config.InitConfig(&cfgFile, &envMode))
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
