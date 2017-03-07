package cmd

import (
	"github.com/abiosoft/ishell"
	cmd_state "github.com/synw/microb-cli/libmicrob/cmd/state"	
)


func GetCmds(shell *ishell.Shell) *ishell.Shell {
	shell.AddCmd(cmd_state.Use())
	return shell
}
