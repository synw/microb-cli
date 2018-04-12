package cmds

import (
	"github.com/synw/microb-cli/libmicrob/msgs"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/types"
)

func GetCmd(cmdName string, cmdArgs []interface{}, state *cliTypes.State) (*types.Cmd, bool) {
	/*
		Get a command from its service
	*/
	var cmdSrv *types.Service
	for sname, srv := range state.Services {
		// check if the first argument is a service name
		if cmdName == sname {
			cmdName = cmdArgs[0].(string)
			if len(cmdArgs) > 1 {
				cmdArgs = cmdArgs[1:]
			}
			cmdSrv = srv
			break
		}
	}
	if cmdSrv == nil {
		if state.CurrentService == nil {
			msgs.Error("No current service set, please set a service: ex: set infos")
			var rcmd *types.Cmd
			return rcmd, false
		}
		cmdSrv = state.CurrentService
	}
	for cname, cmd := range cmdSrv.Cmds {
		if cmdName == cname {
			cmd.Args = cmdArgs
			cmd.From = "cli"
			return cmd, true
		}
	}
	var ecmd *types.Cmd
	return ecmd, false
}

func EnsureCmdsHaveService(state *cliTypes.State) {
	/*
		Make sure all commands have their Service value set
	*/
	for sname, srv := range state.Services {
		for _, cmd := range srv.Cmds {
			cmd.Service = sname
		}
	}
}
