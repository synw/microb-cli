package cmd

import (
	"errors"
	//"github.com/synw/microb-cli/libmicrob/msgs"
	"github.com/synw/microb-cli/libmicrob/state"
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

/*func GetCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["use"] = Use()
	cmds["using"] = Using()
	cmds["ping"] = Ping()
	cmds["services"] = Services()
	return cmds
}*/

func IsValid(cmd *types.Cmd) bool {
	for _, scmd := range state.Cmds {
		if cmd.Name == scmd.Name {
			return true
		}
	}
	return false
}

func serverExists(server_name string) *terr.Trace {
	for name, _ := range state.Servers {
		if server_name == name {
			return nil
		}
	}
	msg := "Server " + server_name + " not found: please check your config file"
	err := errors.New(msg)
	tr := terr.New("ws.serverExists", err)
	return tr
}
