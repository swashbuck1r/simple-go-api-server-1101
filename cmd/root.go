package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/swashbuck1r/simple-go-api-server/server"
)

const (
	APP_ENV_PREFIX = ""
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

			var config server.Configuration
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
	if err := viper.ReadConfig(bytes.NewBuffer(server.DefaultConfiguration)); err != nil {
		log.Fatalln("Failed to read default config file", err)
	}

	cobra.OnInitialize(initConfig)

	// config file flag
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")

	// env mode flag
	rootCmd.PersistentFlags().StringVar(&envMode, "env", "dev", "environment execution mode (dev, production, etc)")
	viper.BindPFlag("env", rootCmd.Flags().Lookup("env"))
}

func initConfig() {
	log.Printf("running in env mode: %s\n", envMode)
	loadEnvDefaults(envMode)

	// allow use of _ as nested property (since . is not supported in env names in most shells)
	viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

	// use env vars with the prefix as config overrides
	if APP_ENV_PREFIX != "" {
		viper.SetEnvPrefix(APP_ENV_PREFIX)
	}

	// use env vars as override in viper config
	viper.AutomaticEnv()

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			log.Fatalln("Failed to read config file:", err)
		}
	}
	viper.Set("env", envMode)
}

func loadEnvDefaults(env string) {
	if "" == env {
		env = "dev"
	}

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}

func PrettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "  ")
	return string(s)
}
