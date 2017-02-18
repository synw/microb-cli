package conf

import (
	"fmt"
	"github.com/spf13/viper"
)


func GetCliConf() map[string]interface{} {
	viper.SetConfigName("dev_config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.microb-cli")
	err := viper.ReadInConfig()
	if err != nil {
	    panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	conf := make(map[string]interface{})
	conf["servers"] = viper.Get("servers")
	viper.SetDefault("debug", false)
	conf["debug"] = viper.Get("debug")
	return conf
}
