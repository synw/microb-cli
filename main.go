package main

import (
	"fmt"
	"errors"
	"flag"
	"github.com/abiosoft/ishell"
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb-cli/libmicrob/metadata"
	//"github.com/synw/microb-cli/libmicrob/listeners"
	"github.com/synw/microb-cli/libmicrob/cmds"
)


var CurrentServer *datatypes.Server

var UseServerFlag = flag.String("s", "", "Use server: -s=server_name")
var Servers map[string]*datatypes.Server = metadata.GetServers()

var shell = ishell.New()

func main(){
	flag.Parse()
	
	// listen to feedback channel and print
	/*
	go func() {
        for msg := range(FeedbackChan) {
            shell.Println(msg)
        }
    }()*/
	
	// INIT
	var msg string
	if *UseServerFlag != "" {
		_, err  := cmds.ServerExists(*UseServerFlag)
		if err != nil {
			shell.Println(err)
			return
		}
		CurrentServer= Servers[*UseServerFlag]
		msg = "Using server "+CurrentServer.Domain
	} else {
		msg = "Available servers: "
		for _, server := range Servers {
			msg = msg+server.Domain+" "
		}
	}
    shell.Println(msg)
    
    // USE
    shell.AddCmd(&ishell.Cmd{
		Name: "use",
		Help: "Use server: use server_domain",
		Func: func(c *ishell.Context) {
			if len(c.Args) == 0 {
				c.Err(errors.New("missing server domain"))
				return
			}
			if len(c.Args) > 1 {
				c.Err(errors.New("please use only one server at the time"))
				return
			}
			server_name := c.Args[0]
			_, err  := cmds.ServerExists(server_name)
			if err != nil {
				fmt.Println(err)
				return
			}
			CurrentServer= Servers[server_name]
			c.Println("Using server", server_name)
		},
	})
	
	// USING
	shell.AddCmd(&ishell.Cmd{
        Name: "using",
        Help: "Check what server is in use",
        Func: func(c *ishell.Context) {
        	if CurrentServer.Domain != "server_not_set" {
            	c.Println("Using server", CurrentServer.Domain)
            } else {
            	c.Println("No server in use. To select a server type: use server_name")
            }
        },
    })
    
    // PING
    shell.AddCmd(&ishell.Cmd{
        Name: "ping",
        Help: "Ping the current server",
        Func: func(ctx *ishell.Context) {
        	_, err, msg := cmds.SendCmd(ctx, "ping", CurrentServer)
        	if err != nil {
        		ctx.Err(err)
        	} else {
        		shell.Println(msg)
        	}
        },
    })   
    shell.Start()
}
