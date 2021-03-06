package conf

import (
	"errors"
	"github.com/spf13/viper"
	globalConf "github.com/synw/microb/conf"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
)

func GetServers(dev_mode bool) (map[string]*types.Server, *terr.Trace) {
	if dev_mode {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.microb-cli")
	err := viper.ReadInConfig()
	servers := make(map[string]*types.Server)
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
		name := sv["name"].(string)
		wsaddr := sv["centrifugo_addr"].(string)
		wskey := sv["centrifugo_key"].(string)
		comchan_in, comchan_out := globalConf.GetComChan(name)
		servers[name] = &types.Server{name, wsaddr, wskey, comchan_in, comchan_out}
	}
	return servers, nil
}
