package state

import (
	"errors"
	"fmt"
	"github.com/synw/centcom"
	"github.com/synw/microb-cli/cmd"
	"github.com/synw/microb-cli/conf"
	"github.com/synw/microb/cmd"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
	"github.com/synw/microb-cli/cmd/info"
)

var Servers map[string]*types.Server
var Server *types.Server
var Cli *centcom.Cli
var Verbosity int = 1
var Cmds []*types.Command

func InitState(dev_mode bool, verbosity int) *terr.Trace {
	Verbosity = verbosity
	servers, trace := conf.GetServers(dev_mode)
	if trace != nil {
		trace := terr.Pass("state.InitState", trace)
		return trace
	}
	Servers = servers
	if Verbosity > 2 {
		msg := "Found servers "
		for name, _ := range Servers {
			msg = msg + name + " "
		}
		fmt.Println(msg)
	}
	Cmds = GetCmds()
	return nil
}

func GetCmds() []*types.Command { 
	ping = info.Ping()
	return [ping]
}

func InitServer() *terr.Trace {
	cli := centcom.NewClient(Server.WsAddr, Server.WsKey)
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
