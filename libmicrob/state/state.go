package state

import (
	"errors"
	"github.com/synw/centcom"
	"github.com/synw/microb-cli/libmicrob/conf"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb-cli/services"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Servers map[string]*types.WsServer
var Server *types.WsServer = nil
var Cli *centcom.Cli
var Services map[string]*types.Service
var Cmds map[string]*types.Cmd

func Init() *terr.Trace {
	// get servers
	servers, srvNames, tr := conf.Get()
	if tr != nil {
		tr := terr.Pass("state.Init", tr)
		return tr
	}
	Servers = servers
	// get services
	Services, tr = services.Get(srvNames)
	if tr != nil {
		tr := terr.Pass("state.Init", tr)
		return tr
	}
	Cmds = make(map[string]*types.Cmd)
	// get all commands
	for _, srv := range Services {
		for cname, cmd := range srv.Cmds {
			Cmds[cname] = cmd
		}
	}
	return nil
}

func InitServer() *terr.Trace {
	cli := centcom.NewClient(Server.Addr, Server.Key)
	err := centcom.Connect(cli)
	if err != nil {
		tr := terr.New("state.InitServer", err)
		return tr
	}
	cli.IsConnected = true
	msg := "Client connected: using command channel " + Server.CmdChanIn
	msgs.Ok(msg)
	err = cli.CheckHttp()
	if err != nil {
		tr := terr.New("state.InitServer", err)
		return tr
	}
	msgs.Ok("Http transport ready")
	Cli = cli
	err = Cli.Subscribe(Server.CmdChanOut)
	if err != nil {
		err := errors.New(err.Error())
		tr := terr.New("state.InitServer", err)
		return tr
	}
	return nil
}
