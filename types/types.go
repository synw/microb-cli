package types

import (
	"github.com/synw/centcom"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
)

type State struct {
	WsServers      map[string]*types.WsServer
	WsServer       *types.WsServer
	Cli            *centcom.Cli
	Services       map[string]*types.Service
	Conf           *types.Conf
	InitServer     func() *terr.Trace
	CurrentService *types.Service
	CmdIds         []string
}
