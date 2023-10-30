package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

var LoginUrl string
var ClassUrl string
var Port string

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalln(err)
	}
	LoginUrl = viper.GetString("login_url")
	ClassUrl = viper.GetString("class_url")
	Port = viper.GetString("port")
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
}
