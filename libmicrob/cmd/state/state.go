package state

import (
	"fmt"
	"github.com/abiosoft/ishell"
	//"github.com/synw/terr"
)


func Use() *ishell.Cmd {
	command := &ishell.Cmd{
        Name: 	"use",
        Help: 	"Use server: use server_domain",
        Func: 	func(ctx *ishell.Context) {
					var err error
					if len(ctx.Args) == 0 {
						err = errors.New("missing server domain")
						return
					}
					if len(ctx.Args) > 1 {
						err = errors.New("please use only one server at the time")
						return
					}
					server_name := ctx.Args[0]
					_, err  = cmds.ServerExists(server_name)
					if err != nil {
						ctx.Println(format.ErrorFormated(err))
						return
					}
					CurrentServer= Servers[server_name]
					ctx.Println("Using server", server_name)
				},
    }
	return command
}
