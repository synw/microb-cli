package conf

import (
	"github.com/spf13/viper"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
)

func getComChan(name string) (string, string) {
	comchan_in := "cmd:$" + name + "_in"
	comchan_out := "cmd:$" + name + "_out"
	return comchan_in, comchan_out
}

func Get() (map[string]*types.WsServer, []string, *terr.Trace) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("~/.microb-cli")
	err := viper.ReadInConfig()
	servers := make(map[string]*types.WsServer)
	var srvs []string
	if err != nil {
		switch err.(type) {
		case viper.ConfigParseError:
			tr := terr.New(err)
			return servers, srvs, tr
		default:
			tr := terr.New("Unable to locate config file")
			return servers, srvs, tr
		}
	}
	available_servers := viper.Get("servers").([]interface{})
	for i, _ := range available_servers {
		sv := available_servers[i].(map[string]interface{})
		name := sv["name"].(string)
		wsaddr := sv["centrifugo_addr"].(string)
		wskey := sv["centrifugo_key"].(string)
		comchan_in, comchan_out := getComChan(name)
		servers[name] = &types.WsServer{name, wsaddr, wskey, comchan_in, comchan_out}
	}
	snames := viper.Get("services").([]interface{})
	for _, srv := range snames {
		srvs = append(srvs, srv.(string))
	}

	return servers, srvs, nil
}
