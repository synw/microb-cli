package conf

import (
	"errors"
	"github.com/spf13/viper"
	"github.com/synw/terr"
	"github.com/synw/microb/libmicrob/datatypes"
)


func GetServers(dev_mode bool) (map[string]*datatypes.Server, *terr.Trace) {
	if dev_mode {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.microb-cli")
	err := viper.ReadInConfig()
	servers := make(map[string]*datatypes.Server)
	if err != nil {
		switch err.(type) {
		case viper.ConfigParseError:
			trace := terr.New("getServers", err)
			return servers, trace
		default:
			err := errors.New("Unable to locate config file")
			trace := terr.New("getServers", err)
			return servers, trace
		}
	}
	available_servers := viper.Get("servers").([]interface{})	
	for i, _ := range available_servers {
		sv := available_servers[i].(map[string]interface{})
		domain := sv["domain"].(string)
		host :=  sv["http_host"].(string)
		port := int(sv["http_port"].(float64))
		wshost := sv["centrifugo_host"].(string)
		wsport := int(sv["centrifugo_port"].(float64))
		wskey := sv["centrifugo_key"].(string)
		comchan := sv["command_channel"].(string)		
		servers[domain] = &datatypes.Server{domain, host, port, wshost, wsport, wskey, comchan}
	}
	return servers, nil
}
