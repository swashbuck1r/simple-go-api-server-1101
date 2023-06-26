package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

const (
	APP_ENV_PREFIX = ""
)

// Initializes the viper configuration based on a config file and env mode
func InitConfig(cfgFile *string, envMode *string) func() {
	return func() {
		log.Printf("running in env mode: %s\n", *envMode)
		loadEnvDefaults(*envMode)

		// allow use of _ as nested property (since . is not supported in env names in most shells)
		viper.SetEnvKeyReplacer(strings.NewReplacer(`.`, `_`))

		// use env vars with the prefix as config overrides
		if APP_ENV_PREFIX != "" {
			viper.SetEnvPrefix(APP_ENV_PREFIX)
		}

		// use env vars as override in viper config
		viper.AutomaticEnv()

		if *cfgFile != "" {
			// Use config file from the flag.
			viper.SetConfigFile(*cfgFile)
			if err := viper.ReadInConfig(); err == nil {
				fmt.Println("Using config file:", viper.ConfigFileUsed())
			} else {
				log.Fatalln("Failed to read config file:", err)
			}
		}
		viper.Set("env", envMode)
	}
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
