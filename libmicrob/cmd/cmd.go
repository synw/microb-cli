package cmd

import (
	"github.com/abiosoft/ishell"
	"github.com/synw/microb-cli/libmicrob/cmd/info"
	cmd_state "github.com/synw/microb-cli/libmicrob/cmd/state"
)

func GetCmds(shell *ishell.Shell) *ishell.Shell {
	shell.AddCmd(cmd_state.Use())
	shell.AddCmd(cmd_state.Using())
	shell.AddCmd(info.Ping())
	return shell
}
