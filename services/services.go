package services

import (
	"github.com/synw/microb/libmicrob/types"
)

func Init() {

}

/*
func New(name string) *types.Service {
	serv := *types.Service{Name: name}
	return serv
}

func GetService(service *types.Service) *types.Service {
	service.Cmds = GetCmds(service)
}

func GetCmds(service *types.Service) map[string]*types.Cmd {
	cmds := service.GetServerCmds()
	cliCmds := service.GetCliCmds()
	for k, v := range cmds {
		cmds[k] = v
	}
	return cmds
}*/
