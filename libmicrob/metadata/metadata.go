package metadata

import (
	//"fmt"
	"errors"
    "github.com/synw/microb/libmicrob/datatypes"
    "github.com/synw/microb-cli/libmicrob/conf"
)


var Config = conf.GetCliConf()

func IsDebug() bool {
	d := Config["debug"].(bool)
	return d
}

func newServer(domain string, host string, port string, ws_host string, ws_port string, ws_key string) *datatypes.Server {
	server := &datatypes.Server{
								Domain:domain, 
								Host:host, 
								Port:port, 
								WebsocketsHost:ws_host, 
								WebsocketsPort:ws_port, 
								WebsocketsKey:ws_key,
								}
	return server
}

func GetServers() map[string]*datatypes.Server {
	available_servers := Config["servers"].([]interface{})
	servers := make(map[string]*datatypes.Server)
	for i, _ := range available_servers {
		sv := available_servers[i].(map[string]interface{})
		domain := sv["domain"].(string)
		host :=  sv["http_host"].(string)
		port := sv["http_port"].(string)
		ws_host := sv["centrifugo_host"].(string)
		ws_port := sv["centrifugo_port"].(string)
		ws_key := sv["centrifugo_secret_key"].(string)
		servers[domain] = newServer(domain, host, port, ws_host, ws_port, ws_key)
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
