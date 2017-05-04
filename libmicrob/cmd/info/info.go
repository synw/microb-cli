package info

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/terr"
	command "github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
)


// PING
func Ping() *ishell.Cmd {
    command := &ishell.Cmd{
        Name: "ping",
        Help: "Ping the current server",
        Func: func(ctx *ishell.Context) {
        	cmd := command.New("ping", "cli", "")
        	cmd, timeout, trace := handler.SendCmd(cmd, ctx)
        	if trace != nil {
        		trace = terr.Pass("cmd.info.Ping", trace)
        		msg := trace.Formatc()
        		ctx.Println(msg)
        		return
        	}
        	if timeout == true {
        		err := terr.Err("Timeout: server is not responding")
				ctx.Println(err.Error())
        	}
        	return
        },
    }
    return command
}

func Http() *ishell.Cmd {
    command := &ishell.Cmd{
        Name: "http",
        Help: "Retrieve the state of the http server",
        Func: func(ctx *ishell.Context) {
        	cmd := command.New("http", "cli", "")
        	cmd, timeout, trace := handler.SendCmd(cmd, ctx)
        	if trace != nil {
        		trace = terr.Pass("cmd.info.Http", trace)
        		msg := trace.Formatc()
        		ctx.Println(msg)
        		return
        	}
        	if timeout == true {
        		err := terr.Err("Timeout: server is not responding")
				ctx.Println(err.Error())
        	}
        	return
        },
    }
    return command
}
