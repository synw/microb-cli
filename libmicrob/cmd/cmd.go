package cmd

import (
	m "github.com/synw/microb/libmicrob"
	"github.com/synw/microb/libmicrob/types"
)

func GetCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["use"] = Use()
	return cmds
}

func Use() *types.Cmd {
	cmd := &types.Cmd{
		Name: "use",
		Exec: use,
	}
	return cmd
}

func use(cmd *types.Cmd) *types.Cmd {
	server := cmd.Args[0].(string)
	m.Debug(server)
	tr := serverExists(server_name)
	if tr != nil {
		tr = terr.Pass("comd.state.Use", tr)
		tr.Error()
		return
	}
	/*state.Server = state.Servers[server_name]
	// init cli and check server
	tr = state.InitServer()
	if tr != nil {
		err := errors.New("can not connect to websockets server: check your config")
		tr := terr.Add("cmd.state.Use", err, tr)

	} else {
		msg := "Using server " + server_name
		ctx.Println(msg)
	}*/
	return cmd
}

func serverExists(server_name string) *terr.trace {
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
