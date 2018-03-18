package types

import (
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

type State struct {
	WsServers  map[string]*types.WsServer
	WsServer   *types.WsServer
	Cli        *centcom.Cli
	Services   map[string]*types.Service
	Cmds       map[string]*types.Cmd
	Conf       *types.Conf
	InitServer func() *terr.Trace
}
