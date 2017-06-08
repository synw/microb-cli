package services

import (
	"github.com/abiosoft/ishell"
)

func GetCmds(shell *ishell.Shell) *ishell.Shell {
	for _, cmd := range commands {
		shell.AddCmd(cmd)
	}
	return shell
}
