package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// InitConfig initializes the configuration for the application.
//
// It reads the value of the "APP_ENV" environment variable and uses it to set the
// configuration file name and path. Then, it reads the configuration file using
// Viper library. If there is an error reading the file, it falls back to reading
// the values from the OS environment. Finally, it logs the values of all the
// environment variables used by the application.
func InitConfig() {
	environment := os.Getenv("APP_ENV")

	log.Printf("environment %s", environment)
	viper.SetConfigName(".env." + environment)
	viper.AddConfigPath("./config/env/")
	viper.SetConfigType("env")

	err := viper.ReadInConfig()
	allEnv := [4]string{
		"DB_URL",
		"DEBUG",
		"PORT",
		"TESTING",
	}
	if err != nil {
		log.Printf("WARN: error loading .env.%s file. Loading from OS environment..", environment)
		for i := range allEnv {
			viper.SetDefault(allEnv[i], os.Getenv(allEnv[i]))
		}
	}
	log.Println("ENVIRONMENTS >>>")
	for i := range allEnv {
		log.Printf("%s : %s", allEnv[i], viper.GetString(allEnv[i]))
	}
}

func GetDbUrl() string {
	return viper.GetString("DB_URL")
}

func GetDebug() bool {
	return viper.GetBool("DEBUG")
}

func GetPort() string {
	return viper.GetString("PORT")
}

func GetTesting() bool {
	return viper.GetBool("TESTING")
}
