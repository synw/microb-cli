package httpServer

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/libmicrob/cmd/handler"
	command "github.com/synw/microb/libmicrob/cmd"
	"github.com/synw/microb/services"
	"github.com/synw/terr"
)

func Start() *ishell.Cmd {
	command := &ishell.Cmd{
		Name: "start",
		Help: "Start the http server",
		Func: func(ctx *ishell.Context) {
			cmd := command.New("start", services.HttpService, "cli", "")
			cmd, timeout, tr := handler.SendCmd(cmd, ctx)
			if tr != nil {
				tr = terr.Pass("cmd.httpServer.Start", tr)
				msg := tr.Formatc()
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
			cmd := command.New("stop", services.HttpService, "http", "cli", "")
			cmd, timeout, tr := handler.SendCmd(cmd, ctx)
			if tr != nil {
				tr = terr.Pass("cmd.httpServer.Stop", tr)
				msg := tr.Formatc()
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
