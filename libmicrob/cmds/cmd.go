package cmds

import (
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

func IsValid(cmd *types.Cmd, state *cliTypes.State) bool {
	for _, scmd := range state.Cmds {
		if cmd.Name == scmd.Name {
			return true
		}
	}
	return false
}
