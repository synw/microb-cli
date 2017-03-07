package state

import (
	"fmt"
	"github.com/synw/terr"
	"github.com/synw/centcom"
	"github.com/synw/centcom/ws"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb-cli/libmicrob/conf"
)


var Servers map[string]*datatypes.Server
var Server *datatypes.Server
var Verbosity int = 1


func InitState(dev_mode bool, verbosity int) (*terr.Trace) {
	Verbosity = verbosity
	servers, trace := conf.GetServers(dev_mode)
	if trace != nil {
		trace := terr.Pass("state.InitState", trace)
		return trace
	}
	Servers = servers
	if Verbosity > 2 {
		msg := "Found servers "
		for name, _ := range(Servers) {
			msg = msg+name+" "
		}
		fmt.Println(msg)
	}
	return nil
}

// internal methods

func InitWsCli() (*ws.Cli, *terr.Trace) {
	cli := ws.NewClient(Server.WsHost, Server.WsPort, Server.WsKey)
	cli, err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("ws.InitCli", err)
		var cli *ws.Cli
		return cli, trace
	}
	cli.IsConnected = true
	if Verbosity > 1 {
		fmt.Println(terr.Ok("Websockets client connected"))
	}
	cli, err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("ws.InitCli", err)
		return cli, trace
	}
	if Verbosity > 1 {
		fmt.Println(terr.Ok("Websockets http transport ready"))
	}	
	return cli, nil
}
