package state

import (
	"github.com/synw/centcom"
	"github.com/synw/microb-cli/libmicrob/cmds"
	"github.com/synw/microb-cli/libmicrob/conf"
	"github.com/synw/microb-cli/libmicrob/msgs"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb-cli/services"
	//"github.com/synw/microb/libmicrob/types"
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
		tr = tr.Add("Can not initialize websockets server", "fatal")
		return State, tr
	}
	// get services
	State.Services, tr = services.GetAll(srvNames)
	if tr != nil {
		tr := tr.Pass()
		return State, tr
	}
	cmds.EnsureCmdsHaveService(State)
	return State, nil
}

func initServer() *terr.Trace {
	cli := centcom.NewClient(State.WsServer.Addr, State.WsServer.Key)
	err := centcom.Connect(cli)
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	cli.IsConnected = true
	msg := "Client connected: using command channel " + State.WsServer.CmdChanIn
	msgs.Ok(msg)
	err = cli.CheckHttp()
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	msgs.Ok("Http transport ready")
	err = cli.Subscribe(State.WsServer.CmdChanOut)
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	State.Cli = cli
	return nil
}
