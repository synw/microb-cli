package services

import (
	color "github.com/acmacalister/skittles"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func Get(servs []string) (map[string]*types.Service, *terr.Trace) {
	srvs := make(map[string]*types.Service)
	manSrvs := services
	for _, name := range servs {
		for k, v := range manSrvs {
			if k == name {
				srvs[k] = v
				msgs.Status("Initializing " + color.BoldWhite(v.Name) + " service")
				break
			}
		}
		/*srv, err := srv.Init(m.Verbose())
		if err != nil {
			tr := terr.New("services.Init", "Can not initialize service "+name, err)
			return srvs, tr
		}*/
	}
	return srvs, nil
}
