package main

import (
	"errors"
	"flag"
	"github.com/synw/microb-cli/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb-cli/libmicrob/prompter"
	"github.com/synw/microb-cli/libmicrob/state"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var server = flag.String("u", "__unset__", "Use server")

func main() {
	flag.Parse()
	// read conf
	tr := state.Init()
	if tr != nil {
		err := errors.New("Unable to init state")
		tr := terr.Add("main", err, tr)
		terr.Fatal("main", tr)
	}
	msgs.Ok("State initialized")
	srvs := "Available servers:"
	for name, _ := range state.Servers {
		srvs = srvs + " " + name
	}
	if *server != "__unset__" {
		com := cmd.Use()
		var args []interface{}
		args = append(args, *server)
		com.Args = args
		_, tr := com.ExecCli.(func(*types.Cmd) (*types.Cmd, *terr.Trace))(com)
		if tr != nil {
			tr.Formatc()
		}
	}
	msgs.Msg(srvs)
	prompter.Prompt()
}
