package state

import (
	"errors"
	"github.com/synw/centcom"
	"github.com/synw/microb-cli/libmicrob/conf"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Servers map[string]*types.WsServer
var Server *types.WsServer = nil
var Cli *centcom.Cli
var Verbosity int = 1
var Services []*types.Service

func Init(verbosity int) *terr.Trace {
	Verbosity = verbosity
	servers, trace := conf.GetServers()
	if trace != nil {
		trace := terr.Pass("state.Init", trace)
		return trace
	}
	Servers = servers
	if Verbosity > 2 {
		msg := "Found servers "
		for name, _ := range Servers {
			msg = msg + name + " "
		}
		msgs.Msg(msg)
	}
	return nil
}

func InitServer() *terr.Trace {
	cli := centcom.NewClient(Server.Addr, Server.Key)
	err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("centcom.InitCli", err)
		return trace
	}
	cli.IsConnected = true
	if Verbosity > 0 {
		msg := "Client connected: using command channel " + Server.CmdChanIn
		terr.Ok(msg)
	}
	err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("centcom.InitCli", err)
		return trace
	}
	if Verbosity > 0 {
		terr.Ok("Http transport ready")
	}
	Cli = cli
	err = Cli.Subscribe(Server.CmdChanOut)
	if err != nil {
		err := errors.New(err.Error())
		trace := terr.New("state.InitServer", err)
		return trace
	}
	return nil
}
