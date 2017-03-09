package state

import (
	"fmt"
	"github.com/synw/terr"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb-cli/libmicrob/conf"
)


var Servers map[string]*datatypes.Server
var Server *datatypes.Server
var Cli *centcom.Cli
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

func InitServer() *terr.Trace {
	centcom.SetVerbosity(Verbosity)
	cli := centcom.NewClient(Server.WsHost, Server.WsPort, Server.WsKey)
	cli, err := centcom.Connect(cli)
	if err != nil {
		trace := terr.New("centcom.InitCli", err)
		return trace
	}
	cli.IsConnected = true
	if Verbosity > 0 {
		msg := "Client connected: using command channel "+Server.CmdChannel
		fmt.Println(terr.Ok(msg))
	}
	cli, err = cli.CheckHttp()
	if err != nil {
		trace := terr.New("centcom.InitCli", err)
		return trace
	}
	if Verbosity > 0 {
		fmt.Println(terr.Ok("Http transport ready"))
	}
	Cli = cli
	return nil
}
