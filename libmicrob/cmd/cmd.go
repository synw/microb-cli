package cmd

import (
	"errors"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb-cli/libmicrob/state"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

func GetCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["use"] = Use()
	cmds["using"] = Using()
	cmds["ping"] = Ping()
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

func Ping() *types.Cmd {
	cmd := &types.Cmd{
		Name:    "ping",
		Service: "infos",
		Exec:    ping,
	}
	return cmd
}

func ping(cmd *types.Cmd) *types.Cmd {
	cmd, timeout, tr := handler.SendCmd(cmd)
	if tr != nil {
		tr = terr.Pass("cmd.Ping", tr)
		tr.Print()
		return cmd
	}
	if timeout == true {
		msgs.Timeout("The server is not responding")
	}
	return cmd
}

func using(cmd *types.Cmd) *types.Cmd {
	if state.Server == nil {
		msgs.Warning("No server selected: try the use command: ex: " + msgs.Bold("use") + " server1")
	} else {
		msgs.Msg("Using server " + msgs.Bold(state.Server.Name))
	}
	return cmd
}

func use(cmd *types.Cmd) *types.Cmd {
	if len(cmd.Args) != 1 {
		msgs.Warning("Please provide a server name: ex: use localhost")
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
		msgs.Ready(msg)
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
