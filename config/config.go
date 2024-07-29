package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)


func GetLocalEnv(key string) string {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName(".env")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Find and read the config file

	if err := viper.ReadInConfig(); err != nil {
		log.Panic(fmt.Printf("Error while reading config file %s", err))
		return ""
	}

	// If the type is a string then ok will be true
	// ok will make sure the program not break
	if value, ok := viper.Get(key).(string); !ok {
		log.Panic("Invalid type assertion", key)
		return ""
	} else {
		return value
	}
}
