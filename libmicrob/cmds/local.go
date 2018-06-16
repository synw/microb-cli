package cmds

import (
	"errors"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func Use() *types.Cmd {
	cmd := &types.Cmd{
		Name:    "use",
		ExecCli: runUse,
	}
	return cmd
}

func Using() *types.Cmd {
	cmd := &types.Cmd{
		Name:    "using",
		ExecCli: runUsing,
	}
	return cmd
}

func runUsing(cmd *types.Cmd, state *cliTypes.State) (*types.Cmd, *terr.Trace) {
	if state.WsServer == nil {
		msgs.Warning("No server selected: try the use command: ex: " + msgs.Bold("use") + " server1")
	} else {
		msgs.Msg("Using server " + msgs.Bold(state.WsServer.Name))
	}
	return cmd, nil
}

func runUse(cmd *types.Cmd, state *cliTypes.State) (*types.Cmd, *terr.Trace) {
	if len(cmd.Args) != 1 {
		msgs.Warning("Please provide a server name: ex: use localhost")
		tr := terr.New("Can not find server")
		return cmd, tr
	}
	server := cmd.Args[0].(string)
	msgs.Status("Connecting to server " + server + " ...")
	tr := serverExists(server, state)
	if tr != nil {
		tr = tr.Pass()
		return cmd, tr
	}
	state.WsServer = state.WsServers[server]
	tr = state.InitServer()
	// init cli and check server
	if tr != nil {
		err := errors.New("can not connect to websockets server: check your config")
		tr := tr.Add(err)
		events.Error("microb", "Can not connect to websockets server", tr)
		return cmd, tr
	} else {
		msg := "Connnected to server " + server
		msgs.Ready(msg)
	}
	return cmd, nil
}

func serverExists(server_name string, state *cliTypes.State) *terr.Trace {
	for name, _ := range state.WsServers {
		if server_name == name {
			return nil
		}
	}
	msg := "Server " + server_name + " not found: please check your config file"
	tr := terr.New(msg)
	return tr
}
