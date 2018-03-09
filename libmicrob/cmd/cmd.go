package cmd

import (
	"errors"
	"github.com/synw/microb-cli/libmicrob/state"
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["use"] = Use()
	cmds["using"] = Using()
	return cmds
}

func Use() *types.Cmd {
	cmd := &types.Cmd{
		Name: "use",
		Exec: use,
	}
	return cmd
}

func Using() *types.Cmd {
	cmd := &types.Cmd{
		Name: "using",
		Exec: using,
	}
	return cmd
}

func using(cmd *types.Cmd) *types.Cmd {
	if state.Server == nil {
		m.Warning("No server selected: try the use command: ex: "+m.Bold("use")+" server1", false)
	} else {
		m.Msg("Using server "+m.Bold(state.Server.Name), false)
	}
	return cmd
}

func use(cmd *types.Cmd) *types.Cmd {
	if len(cmd.Args) != 1 {
		m.Warning("Please provide a server name: ex: use localhost", true)
		return cmd
	}
	server := cmd.Args[0].(string)
	tr := serverExists(server)
	if tr != nil {
		tr = terr.Pass("comd.state.Use", tr)
		tr.Printc()
		return cmd
	}
	state.Server = state.Servers[server]
	tr = state.InitServer()
	if tr != nil {
		tr = terr.Pass("comd.state.Use", tr)
		tr.Printc()
		return cmd
	}
	// init cli and check server
	if tr != nil {
		err := errors.New("can not connect to websockets server: check your config")
		tr := terr.Add("cmd.state.Use", err, tr)
		tr.Printc()
		return cmd
	} else {
		msg := "Connnected to server " + server
		m.Ready(msg, true)
	}
	return cmd
}

func serverExists(server_name string) *terr.Trace {
	for name, _ := range state.Servers {
		if server_name == name {
			return nil
		}
	}
	msg := "Server " + server_name + " not found: please check your config file"
	err := errors.New(msg)
	tr := terr.New("ws.serverExists", err)
	return tr
}
