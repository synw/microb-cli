package main

import (
	"flag"
	"github.com/synw/microb-cli/cmds"
	"github.com/synw/microb-cli/msgs"
	"github.com/synw/microb-cli/prompter"
	st "github.com/synw/microb-cli/state"
	cliTypes "github.com/synw/microb-cli/types"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
)

var server = flag.String("u", "__unset__", "Use server")

func main() {
	flag.Parse()
	// read conf
	state, tr := st.Init()
	if tr != nil {
		tr := tr.Pass()
		tr.Fatal("Unable to init state")
	}
	msgs.Ok("State initialized")
	srvs := "Available servers:"
	for name, _ := range state.WsServers {
		srvs = srvs + " " + name
	}
	// set the server to use if requested
	if *server != "__unset__" {
		cmd := cmds.Use()
		var args []interface{}
		args = append(args, *server)
		cmd.Args = args
		_, tr := cmd.ExecCli.(func(*types.Cmd, *cliTypes.State) (*types.Cmd, *terr.Trace))(cmd, state)
		if tr != nil {
			tr.Print()
		}
	}
	msgs.Msg(srvs)
	prompter.Prompt()
}
