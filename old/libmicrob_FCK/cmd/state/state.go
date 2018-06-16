package state

import (
	//"fmt"
	"errors"
	"github.com/abiosoft/ishell"
	"github.com/acmacalister/skittles"
	"github.com/synw/microb-cli/state"
	"github.com/synw/terr"
)

func Using() *ishell.Cmd {
	command := &ishell.Cmd{
		Name: "using",
		Help: "Server actually in use",
		Func: func(ctx *ishell.Context) {
			if state.Server == nil {
				ctx.Println("No server selected: try the use command: ex:", skittles.BoldWhite("use"), "server1")
			} else {
				ctx.Println("Using server", state.Server.Name)
			}
		},
	}
	return command
}

func Use() *ishell.Cmd {
	command := &ishell.Cmd{
		Name: "use",
		Help: "Use server: use server_domain",
		Func: func(ctx *ishell.Context) {
			if len(ctx.Args) == 0 {
				err := terr.Err("missing server domain")
				ctx.Println(err.Error())
				return
			}
			if len(ctx.Args) > 1 {
				err := terr.Err("please use only one server at the time")
				ctx.Println(err.Error())
				return
			}
			server_name := ctx.Args[0]
			trace := serverExists(server_name)
			if trace != nil {
				trace = terr.Pass("comd.state.Use", trace)
				ctx.Println(trace.Formatc())
				return
			}
			state.Server = state.Servers[server_name]
			// init cli and check server
			trace = state.InitServer()
			if trace != nil {
				err := errors.New("can not connect to websockets server: check your config")
				trace := terr.Add("cmd.state.Use", err, trace)
				ctx.Println(trace.Formatc())
			} else {
				msg := "Using server " + server_name
				ctx.Println(msg)
			}
		},
	}
	return command
}

func serverExists(server_name string) *terr.Trace {
	for name, _ := range state.Servers {
		if server_name == name {
			return nil
		}
	}
	msg := "Server " + server_name + " not found: please check your config file"
	err := errors.New(msg)
	trace := terr.New("ws.serverExists", err)
	return trace
}
