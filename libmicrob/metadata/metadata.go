package metadata

import (
	//"fmt"
	"errors"
	"github.com/synw/microb-cli/libmicrob/conf"
    "github.com/synw/microb/libmicrob/datatypes"
)


var Config = conf.GetCliConf()

func IsDebug() bool {
	d := Config["debug"].(bool)
	return d
}

func GetServers() map[string]*datatypes.Server {
	available_servers := Config["servers"].([]interface{})
	servers := make(map[string]*datatypes.Server)
	for i, _ := range available_servers {
		sv := available_servers[i].(map[string]interface{})
		domain := sv["domain"].(string)
		http_host :=  sv["http_host"].(string)
		http_port := sv["http_port"].(string)
		websockets_host := sv["centrifugo_host"].(string)
		websockets_port := sv["centrifugo_port"].(string)
		websockets_key := sv["centrifugo_secret_key"].(string)
		servers[domain] = &datatypes.Server{domain, http_host, http_port, websockets_host, websockets_port, websockets_key}
	}
	return servers
}

func GetServer(domain string) (*datatypes.Server, error) {
	servers := GetServers()
	for _, server := range(servers) {
		if server.Domain == domain {
			return server, nil
		}
	}
	var s *datatypes.Server
	msg := "Can not find server "+domain
	err := errors.New(msg)
	return s, err
}
