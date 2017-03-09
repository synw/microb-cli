package info

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/terr"	
	"github.com/synw/microb-cli/libmicrob/cmd/cmdhandler"
)


// PING
func Ping() *ishell.Cmd {
    command := &ishell.Cmd{
        Name: "ping",
        Help: "Ping the current server",
        Func: func(ctx *ishell.Context) {
        	cmd := cmdhandler.New("ping", "cli", "")
        	_, trace := cmdhandler.SendCmd(cmd)
        	if trace != nil {
        		trace = terr.Pass("cmd.info.Ping", trace)
        		msg := trace.Formatc()
        		ctx.Println(msg)
        	} else {
        		ctx.Println("ok")
        	}
        	return
        },
    }
    return command
}
