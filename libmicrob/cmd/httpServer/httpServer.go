package httpServer

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/terr"
	command "github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
)


func Start() *ishell.Cmd {
    command := &ishell.Cmd{
        Name: "start",
        Help: "Start the http server",
        Func: func(ctx *ishell.Context) {
        	cmd := command.New("start", "cli", "")
        	cmd, timeout, trace := handler.SendCmd(cmd, ctx)
        	if trace != nil {
        		trace = terr.Pass("cmd.httpServer.Start", trace)
        		msg := trace.Formatc()
        		ctx.Println(msg)
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

func Stop() *ishell.Cmd {
    command := &ishell.Cmd{
        Name: "stop",
        Help: "Stop the http server",
        Func: func(ctx *ishell.Context) {
        	cmd := command.New("stop", "cli", "")
        	cmd, timeout, trace := handler.SendCmd(cmd, ctx)
        	if trace != nil {
        		trace = terr.Pass("cmd.httpServer.Stop", trace)
        		msg := trace.Formatc()
        		ctx.Println(msg)
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
