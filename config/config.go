package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Configuration for mongo environment
type Configuration struct {
	Environment string
	Mongo       MongoConfiguration
}

// MongoConfiguration for database server and collection names
type MongoConfiguration struct {
	Server                  string
	Database                string
	MerchantsCollection     string
	MembersCollection       string
	MerchantsTestCollection string
	MembersTestCollection   string
}

// GetConfig setup configuration based on config.yml
func GetConfig() Configuration {
	conf := Configuration{}

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&conf)
	if err != nil {
		panic(err)
	}

	fmt.Println("config", conf)

	return conf
}
