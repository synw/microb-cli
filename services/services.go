package services

import (
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb-cli/libmicrob/msgs"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetAll(servs []string) (map[string]*types.Service, *terr.Trace) {
	srvs := make(map[string]*types.Service)
	for _, name := range servs {
		for k, v := range services {
			if k == name {
				srvs[k] = v
				msgs.Status("Initializing " + color.BoldWhite(v.Name) + " service")
				break
			}
		}
	}
	return srvs, nil
}

func Get(sname string, state *cliTypes.State) (*types.Service, bool) {
	for name, srv := range state.Services {
		if sname == name {
			return srv, true
		}
	}
	var serv *types.Service
	return serv, false
}
