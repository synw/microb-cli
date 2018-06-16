package info

import (
	"github.com/synw/microb-cli/cmd/handler"
	command "github.com/synw/microb/cmd"
	"github.com/synw/microb/types"	
	"github.com/synw/terr"
)

// PING
func Ping() *types.Command {
	cmd := command.New("ping", "info", "cli", "")
	cmd, timeout, trace := handler.SendCmd(cmd, ctx)
	if trace != nil {
		trace = terr.Pass("cmd.info.Ping", trace)
		msg := trace.Formatc()
		fmt.Println(msg)
		return
	}
	if timeout == true {
		err := terr.Err("Timeout: server is not responding")
		fmt.Println(err.Error())
	}
	return cmd
}
