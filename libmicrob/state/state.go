package state

import (
	"errors"
	"github.com/synw/centcom"
	"github.com/synw/microb-cli/libmicrob/cmds"
	"github.com/synw/microb-cli/libmicrob/conf"
	"github.com/synw/microb-cli/libmicrob/msgs"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb-cli/services"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var State *cliTypes.State = &cliTypes.State{}

func Init() (*cliTypes.State, *terr.Trace) {
	// get servers
	var srvNames []string
	var tr *terr.Trace
	State.WsServers, srvNames, tr = conf.Get()
	State.InitServer = initServer
	if tr != nil {
		tr := terr.Pass("State.Init", tr)
		return State, tr
	}
	// get services
	State.Services, tr = services.Get(srvNames)
	if tr != nil {
		tr := terr.Pass("State.Init", tr)
		return State, tr
	}
	State.Cmds = make(map[string]*types.Cmd)
	State.Cmds["using"] = cmds.Using()
	State.Cmds["use"] = cmds.Use()
	// get all commands
	for _, srv := range State.Services {
		for cname, cmd := range srv.Cmds {
			State.Cmds[cname] = cmd
		}
	}
	return State, nil
}

func initServer() *terr.Trace {
	cli := centcom.NewClient(State.WsServer.Addr, State.WsServer.Key)
	err := centcom.Connect(cli)
	if err != nil {
		tr := terr.New("State.InitServer", err)
		return tr
	}
	cli.IsConnected = true
	msg := "Client connected: using command channel " + State.WsServer.CmdChanIn
	msgs.Ok(msg)
	err = cli.CheckHttp()
	if err != nil {
		tr := terr.New("State.InitServer", err)
		return tr
	}
	msgs.Ok("Http transport ready")
	err = cli.Subscribe(State.WsServer.CmdChanOut)
	if err != nil {
		err := errors.New(err.Error())
		tr := terr.New("State.InitServer", err)
		return tr
	}
	State.Cli = cli
	return nil
}
