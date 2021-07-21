package config

import (
	"log"

	"github.com/spf13/viper"
)

type Settings struct {
	DbHost string `mapstructure:"DB_HOST"`
	DbPort string `mapstructure:"DB_PORT"`
	DbUser string `mapstructure:"DB_USER"`
	DbName string `mapstructure:"DB_NAME"`
	Env string `mapstructure:"ENV"`
	JwtSecret string `mapstructure:"JWT_SECRET"`
	JwtExpires string `mapstructure:"JWT_EXPIRES"`
}

func New() *Settings {
	var cfg Settings
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err:=viper.ReadInConfig()
	if err != nil {
		log.Println("No env file found.", err)
	}

	//try to assign read variables into golang struct
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatal("Error trying to unmarshal configuration.", err)
	}

	return &cfg
}