package services

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/services/httpServer"
)

func GetCmds(shell *ishell.Shell) *ishell.Shell {
	shell.AddCmd(httpServer.Start())
	shell.AddCmd(httpServer.Stop())
	shell.AddCmd(httpServer.Http())
	return shell
}
