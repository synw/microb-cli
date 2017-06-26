package services

import (
	"github.com/abiosoft/ishell"
)

func GetCmds(shell *ishell.Shell, dev bool) *ishell.Shell {
	cmds := commands
	if dev == true {
		cmds = commandsDev
	}
	for _, cmd := range cmds {
		shell.AddCmd(cmd)
	}
	return shell
}
