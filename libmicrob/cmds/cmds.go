package cmds

import (
	"github.com/synw/microb-cli/libmicrob/msgs"
	cliTypes "github.com/synw/microb-cli/libmicrob/types"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
	"time"
)

func New(name string, args ...map[string]interface{}) *types.Cmd {
	date := time.Now()
	var service string
	var params []interface{}
	var from string
	var status string
	var errMsg string
	var trace *terr.Trace
	var returnValues []interface{}
	var execCli interface{}
	var exec interface{}
	var execAfter interface{}
	if len(args) == 1 {
		for k, v := range args[0] {
			if k == "service" {
				service = v.(string)
			} else if k == "params" {
				params = v.([]interface{})
			} else if k == "from" {
				from = v.(string)
			} else if k == "status" {
				status = v.(string)
			} else if k == "errMsg" {
				errMsg = v.(string)
			} else if k == "trace" {
				trace = v.(*terr.Trace)
			} else if k == "returnValues" {
				returnValues = v.([]interface{})
			} else if k == "execCli" {
				execCli = v.(interface{})
			} else if k == "exec" {
				exec = v.(interface{})
			} else if k == "execAfter" {
				execAfter = v.(interface{})
			}
		}
	}
	cmd := &types.Cmd{
		Name:         name,
		Date:         date,
		Service:      service,
		Args:         params,
		From:         from,
		Status:       status,
		ErrMsg:       errMsg,
		Trace:        trace,
		ReturnValues: returnValues,
		ExecCli:      execCli,
		Exec:         exec,
		ExecAfter:    execAfter,
	}
	return cmd
}

func GetCmd(cmdName string, cmdArgs []interface{}, state *cliTypes.State) (*types.Cmd, bool) {
	/*
		Get a command from its service
	*/
	var cmdSrv *types.Service
	for sname, srv := range state.Services {
		//msgs.Debug("CMD "+cmdName, "SRV "+sname)

		// check if the first argument is a service name
		if cmdName == sname {
			cmdName = cmdArgs[0].(string)
			if len(cmdArgs) > 1 {
				cmdArgs = cmdArgs[1:]
			}
			cmdSrv = srv
			break
		} else {
			if state.CurrentService == nil {
				msgs.Error("No current service set, please set a service: ex: set infos")
				var rcmd *types.Cmd
				return rcmd, false
			}
			cmdSrv = state.CurrentService
			msgs.Msg("Using current service " + cmdSrv.Name)
		}
	}
	for cname, cmd := range cmdSrv.Cmds {
		if cmdName == cname {
			cmd.Args = cmdArgs
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
