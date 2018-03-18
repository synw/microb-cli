package main

import (
	"errors"
	"flag"
	"github.com/synw/microb-cli/libmicrob/cmds"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb-cli/libmicrob/prompter"
	st "github.com/synw/microb-cli/libmicrob/state"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var server = flag.String("u", "__unset__", "Use server")

func main() {
	flag.Parse()
	// read conf
	state, tr := st.Init()
	if tr != nil {
		err := errors.New("Unable to init state")
		tr := terr.Add("main", err, tr)
		terr.Fatal("main", tr)
	}
	msgs.Ok("State initialized")
	srvs := "Available servers:"
	for name, _ := range state.WsServers {
		srvs = srvs + " " + name
	}
	if *server != "__unset__" {
		cmd := cmds.Use()
		var args []interface{}
		args = append(args, *server)
		cmd.Args = args
		_, tr := cmd.ExecCli.(func(*types.Cmd, *cliTypes.State) (*types.Cmd, *terr.Trace))(cmd, state)
		if tr != nil {
			tr.Formatc()
		}
	}
	msgs.Msg(srvs)
	prompter.Prompt()
}
